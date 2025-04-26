package mongodb

import (
	"context"
	"errors"
	"fmt"
	logger "go-source/pkg/log"
	"go-source/pkg/metric"
	"go-source/pkg/utils"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrContextNotFoundKeyRegion = errors.New("mongo multi conn: context not found key region")
	ErrNotFoundRegion           = errors.New("mongo multi conn: mapping collections not found region")
)

type ModelInterface interface {
	CollectionName() string
	IndexModels() []mongo.IndexModel
}

type Repository[T ModelInterface] struct {
	*mongo.Collection
	*FilterPlayer
	mappingCollections map[string]*mongo.Collection
	err                error
	keyEncrypt         string
	fieldsNameEnc      map[string]bool
}

func NewRepository[T ModelInterface](dbStorage *DatabaseStorage, opts ...*options.CollectionOptions) *Repository[T] {
	log := logger.GetLogger()

	keyEncrypt := os.Getenv(utils.VGREncryptKey)

	var t T
	fieldsNameEnc := readTagEncrypt(t)
	collectionName := t.CollectionName()
	indexModels := t.IndexModels()

	if dbStorage.db != nil {
		collection, err := newRepository(dbStorage.db, collectionName, indexModels, opts...)
		if err != nil {
			log.Fatal().Msgf("new repository error: %v", err)
		}

		return &Repository[T]{
			Collection:    collection,
			keyEncrypt:    keyEncrypt,
			fieldsNameEnc: fieldsNameEnc,
		}
	}

	connNames := GetMappingRepositoryRegion(t.CollectionName())

	if len(connNames) == 0 {
		log.Fatal().Msgf("mongo multi conn: connNames not found for collectionName=%s", collectionName)
	}

	mappingDatabases := make(map[string]*mongo.Database)
	mappingNames := make(map[string]string)

	for _, connName := range connNames {
		splitConnName := strings.Split(connName, "::")
		if len(splitConnName) != 2 {
			log.Fatal().Msgf("mongo multi conn: connName=%s invalid", connName)
		}

		db, ok := dbStorage.mappingDB[connName]
		if !ok {
			log.Fatal().Msgf("mongo multi conn: connName=%s not found in mappingDB", connName)
		}

		region := splitConnName[0]
		name := splitConnName[1]

		mappingDatabases[region] = db
		mappingNames[region] = name
	}

	mappingCollections := make(map[string]*mongo.Collection)

	for region, db := range mappingDatabases {
		collection, err := newRepository(db, collectionName, indexModels, opts...)
		if err != nil {
			connName := fmt.Sprintf("%s::%s", region, mappingNames[region])
			log.Fatal().Msgf("mongo multi conn - connName=%s: new repository error: %v", connName, err)
		}

		mappingCollections[region] = collection
	}

	return &Repository[T]{
		mappingCollections: mappingCollections,
		keyEncrypt:         keyEncrypt,
		fieldsNameEnc:      fieldsNameEnc,
	}
}

func newRepository(db *mongo.Database, collectionName string, indexModels []mongo.IndexModel, opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	log := logger.GetLogger()

	collection := db.Collection(collectionName, opts...)
	if len(indexModels) > 0 {
		go func() {
			_, err := collection.Indexes().CreateMany(context.Background(), indexModels)
			if err != nil {
				log.Error().Msgf("create index collectionName=%v error: %v", collectionName, err)
			}
		}()
	}

	return collection, nil
}

func (r *Repository[T]) NewFilterPlayer(opts ...FilterPlayerOption) *Repository[T] {
	opts = append(opts, WithMetricComponent(r.Collection.Name()))
	filterPlayer := NewFilterPlayer(opts...)

	return &Repository[T]{
		Collection:    r.Collection,
		FilterPlayer:  filterPlayer,
		err:           r.err,
		keyEncrypt:    r.keyEncrypt,
		fieldsNameEnc: r.fieldsNameEnc,
	}
}

func (r *Repository[T]) NewFilterPlayerMultiConn(ctx context.Context) *Repository[T] {
	if r.mappingCollections == nil {
		return nil
	}

	country, ok := ctx.Value(utils.KeyRegion).(string)
	if !ok {
		return &Repository[T]{
			err: ErrContextNotFoundKeyRegion,
		}
	}
	region := GetRegionCountry(country)

	collection, ok := r.mappingCollections[region]
	if !ok {
		return &Repository[T]{
			err: ErrNotFoundRegion,
		}
	}

	return &Repository[T]{
		Collection:    collection,
		FilterPlayer:  NewFilterPlayer(),
		err:           r.err,
		keyEncrypt:    r.keyEncrypt,
		fieldsNameEnc: r.fieldsNameEnc,
	}
}

func (r *Repository[T]) getIndexesName() [][]string {
	log := logger.GetLogger()
	ctx := context.Background()

	// Get indexes
	cursor, err := r.Indexes().List(ctx)
	if err != nil {
		log.Error().Msgf("get indexes error: %v", err)
	}
	defer cursor.Close(ctx)

	var indexes []bson.D
	if err = cursor.All(ctx, &indexes); err != nil {
		log.Error().Msgf("decode indexes error: %v", err)
	}

	var indexesName [][]string
	for _, index := range indexes {
		var keys []string
		for _, item := range index[1].Value.(bson.D) {
			if item.Key == "_id" {
				break
			}
			keys = append(keys, item.Key)
		}

		if len(keys) > 0 {
			indexesName = append(indexesName, keys)
		}
	}

	return indexesName
}

func (r *Repository[T]) checkIndexOfQuery() {
	if len(r.filter) == 0 {
		return
	}

	var keys []string
	for _, item := range r.filter {
		keys = append(keys, item.Key)
	}

	keyFilter := strings.Join(keys, "|")

	log := logger.GetLogger().With().Str("collectionName", r.Collection.Name()).Str("filter", keyFilter).Logger()

	indexesName := r.getIndexesName()

	if len(indexesName) == 0 {
		log.Warn().Msg("database indexes warning: indexes empty")
		return
	}

	sMax := 0
	for _, index := range indexesName {
		s := calcIndex(index, keys)
		if s > sMax {
			sMax = s
		}
	}

	var indexes []string
	for _, index := range indexesName {
		indexes = append(indexes, strings.Join(index, "|"))
	}

	log = log.With().Interface("indexes", indexes).Logger()

	if sMax == 0 {
		log.Warn().Msg("database indexes warning: not found index")
		return
	}

	if sMax < len(keys) {
		log.Warn().Msg("database indexes warning: found index, but may not be effective")
	}
}

func (r *Repository[T]) FindOneDoc(ctx context.Context, opts ...*options.FindOneOptions) (result *T, err error) {
	if r.metricMethod == "" {
		return r.findOneDoc(ctx, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.findOneDoc(ctx, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) findOneDoc(ctx context.Context, opts ...*options.FindOneOptions) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: FindOneDoc")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	opt := r.optsFindOne
	if len(r.sortOne) > 0 {
		opt.Sort = r.sortOne
	}

	opts = append(opts, &opt)

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	var m T
	err := r.Collection.FindOne(ctx, r.filter, opts...).Decode(&m)
	if err != nil {
		return nil, err
	}

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: FindOneDoc.FindOne")
	}

	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return &m, nil
	}

	result, err := utils.StructDecryptTag(m, r.keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (r *Repository[T]) FindDocs(ctx context.Context, opts ...*options.FindOptions) (result []*T, err error) {
	if r.metricMethod == "" {
		return r.findDocs(ctx, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.findDocs(ctx, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) findDocs(ctx context.Context, opts ...*options.FindOptions) ([]*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: FindDocs")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	opt := r.optsFind
	if len(r.sort) > 0 {
		opt.Sort = r.sort
	}

	opts = append(opts, &opt)

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	cs, err := r.Collection.Find(ctx, r.filter, opts...)
	if err != nil {
		return nil, err
	}
	ms := make([]*T, 0)
	err = cs.All(ctx, &ms)
	if err != nil {
		return nil, err
	}

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: FindDocs.Find")
	}

	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return ms, nil
	}

	var sem chan struct{}
	if len(ms) < 200 {
		sem = make(chan struct{}, 1)
	} else {
		sem = make(chan struct{}, 10)
	}

	defer close(sem)

	err = r.decryptDocsEfficiency(ms, sem)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (r *Repository[T]) decryptDocsEfficiency(docs []*T, sem chan struct{}) error {
	var wg sync.WaitGroup
	var err error

	for i, doc := range docs {
		wg.Add(1)
		sem <- struct{}{} // Add to semaphore to control the limit of concurrent tasks

		go func(i int, input *T) {
			defer func() { <-sem }() // Release semaphore when task is done
			defer wg.Done()

			res, errD := utils.StructDecryptTag(*input, r.keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
			if errD != nil {
				err = errD
				return
			}

			docs[i] = &res
		}(i, doc)
	}

	// Wait for all tasks to complete
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[T]) CreateOneDocument(ctx context.Context, document *T) (result *T, err error) {
	if r.metricMethod == "" {
		return r.createOneDocument(ctx, document)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.createOneDocument(ctx, document)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) createOneDocument(ctx context.Context, document *T) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: CreateOneDocument")
		}()
	}

	data := *document
	var err error
	if r.keyEncrypt != "" && len(r.fieldsNameEnc) > 0 {
		data, err = utils.StructEncryptTag(data, r.keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
		if err != nil {
			return nil, err
		}
	}

	t := time.Now()
	doc, err := r.convertToBson(&data)
	if err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	doc["created_at"] = &t
	doc["updated_at"] = &t
	result, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: CreateOneDocument.InsertOne")
	}

	doc["_id"] = result.InsertedID
	entity, _ := r.convertToObject(doc)
	return entity, nil
}

func (r *Repository[T]) CreateManyDocs(ctx context.Context, documents []*T) (result []*T, err error) {
	if r.metricMethod == "" {
		return r.createManyDocs(ctx, documents)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.createManyDocs(ctx, documents)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) createManyDocs(ctx context.Context, documents []*T) ([]*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	t := time.Now()
	var docsProcessed []interface{}
	for _, document := range documents {
		data := *document
		var err error
		if r.keyEncrypt != "" && len(r.fieldsNameEnc) > 0 {
			data, err = utils.StructEncryptTag(data, r.keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
			if err != nil {
				return nil, err
			}
		}

		docP, err := r.convertToBson(&data)
		if err != nil {
			return nil, err
		}

		docP["created_at"] = &t
		docP["updated_at"] = &t
		docsProcessed = append(docsProcessed, docP)
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	result, err := r.Collection.InsertMany(ctx, docsProcessed)
	if err != nil {
		return nil, err
	}

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: CreateManyDocs.InsertMany")
	}

	var entities []*T
	for i, doc := range docsProcessed {
		doc.(bson.M)["_id"] = result.InsertedIDs[i]
		entity, _ := r.convertToObject(doc.(bson.M))
		entities = append(entities, entity)
	}

	return entities, nil
}

// convertToObject convert bson to object
func (r *Repository[T]) convertToObject(b bson.M) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	if b == nil {
		return nil, nil
	}

	bytes, err := bson.Marshal(b)
	if err != nil {
		return nil, err
	}

	var t T
	err = bson.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// convertToBson Go object to map (to get / query)
func (r *Repository[T]) convertToBson(ent *T) (bson.M, error) {
	if r.err != nil {
		return nil, r.err
	}

	if ent == nil {
		return bson.M{}, nil
	}

	sel, err := bson.Marshal(ent)
	if err != nil {
		return nil, err
	}

	obj := bson.M{}
	err = bson.Unmarshal(sel, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Repository[T]) UpdateOneDoc(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	if r.metricMethod == "" {
		return r.updateOneDoc(ctx, update, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.updateOneDoc(ctx, update, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) updateOneDoc(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: UpdateOneDoc")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return r.Collection.UpdateOne(ctx, r.filter, update, opts...)
	}

	updateEnc, err := encryptBsonUpdate(update, r.fieldsNameEnc, r.keyEncrypt)
	if err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.UpdateOne(ctx, r.filter, updateEnc, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: UpdateOneDoc.UpdateOne")
	}

	return rs, err
}

func (r *Repository[T]) UpsertDoc(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	if r.metricMethod == "" {
		return r.upsertDoc(ctx, update, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.upsertDoc(ctx, update, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) upsertDoc(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: UpsertDoc")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	optUpsert := options.Update().SetUpsert(true)
	opts = append(opts, optUpsert)

	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return r.Collection.UpdateOne(ctx, r.filter, update, opts...)
	}

	updateEnc, err := encryptBsonUpdate(update, r.fieldsNameEnc, r.keyEncrypt)
	if err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.UpdateOne(ctx, r.filter, updateEnc, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: UpsertDoc.UpdateOne")
	}

	return rs, err
}

func (r *Repository[T]) UpdateManyDocs(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	if r.metricMethod == "" {
		return r.updateManyDocs(ctx, update, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.updateManyDocs(ctx, update, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) updateManyDocs(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: UpdateManyDocs")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return r.Collection.UpdateMany(ctx, r.filter, update, opts...)
	}

	updateEnc, err := encryptBsonUpdate(update, r.fieldsNameEnc, r.keyEncrypt)
	if err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.UpdateMany(ctx, r.filter, updateEnc, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: UpdateManyDocs.UpdateMany")
	}

	return rs, err
}

func (r *Repository[T]) FindOneAndUpdateDoc(ctx context.Context, update interface{}, opts ...*options.FindOneAndUpdateOptions) (result *T, err error) {
	if r.metricMethod == "" {
		return r.findOneAndUpdateDoc(ctx, update, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.findOneAndUpdateDoc(ctx, update, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) findOneAndUpdateDoc(ctx context.Context, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: FindOneAndUpdateDoc")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	_update := update

	if r.keyEncrypt != "" && len(r.fieldsNameEnc) > 0 {
		updateEnc, err := encryptBsonUpdate(update, r.fieldsNameEnc, r.keyEncrypt)
		if err != nil {
			return nil, err
		}
		_update = updateEnc
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	res := r.Collection.FindOneAndUpdate(ctx, r.filter, _update, opts...)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var m T
	err := res.Decode(&m)
	if err != nil {
		return nil, err
	}

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: FindOneAndUpdateDoc.FindOneAndUpdate")
	}

	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return &m, nil
	}

	result, err := utils.StructDecryptTag(m, r.keyEncrypt, utils.TagNameEncrypt, utils.TagValEncrypt)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository[T]) CountDocs(ctx context.Context, opts ...*options.CountOptions) (result int64, err error) {
	if r.metricMethod == "" {
		return r.countDocs(ctx, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.countDocs(ctx, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) countDocs(ctx context.Context, opts ...*options.CountOptions) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: CountDocs")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return 0, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.CountDocuments(ctx, r.filter, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: CountDocs.CountDocuments")
	}

	return rs, err
}

func (r *Repository[T]) DeleteOneDoc(ctx context.Context, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	if r.metricMethod == "" {
		return r.deleteOneDoc(ctx, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.deleteOneDoc(ctx, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) deleteOneDoc(ctx context.Context, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: DeleteOneDoc")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.DeleteOne(ctx, r.filter, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: DeleteOneDoc.DeleteOne")
	}

	return rs, err
}

func (r *Repository[T]) DeleteManyDocs(ctx context.Context, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	if r.metricMethod == "" {
		return r.deleteManyDocs(ctx, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.deleteManyDocs(ctx, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) deleteManyDocs(ctx context.Context, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: DeleteManyDocs")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.DeleteMany(ctx, r.filter, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: DeleteManyDocs.DeleteMany")
	}

	return rs, err
}

func (r *Repository[T]) DistinctDocs(ctx context.Context, fieldName string, opts ...*options.DistinctOptions) (result []interface{}, err error) {
	if r.metricMethod == "" {
		return r.distinctDocs(ctx, fieldName, opts...)
	}

	_ = metric.NewMongoDBHistogramWithFunc(
		r.metricComponent,
		r.metricMethod,
		func() error {
			result, err = r.distinctDocs(ctx, fieldName, opts...)
			if err != nil {
				return metric.DefaultErr
			}
			return nil
		},
	)
	return
}

func (r *Repository[T]) distinctDocs(ctx context.Context, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Check query index usage
	go r.checkIndexOfQuery()

	// Measure latency
	if shouldMeasureLatency {
		start := time.Now()
		defer func() {
			logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(start)).Msg("mongodb_latency: DistinctDocs")
		}()
	}

	if err := r.filterEncrypt(); err != nil {
		return nil, err
	}

	var startR *time.Time
	if shouldMeasureLatency {
		now := time.Now()
		startR = &now
	}

	rs, err := r.Collection.Distinct(ctx, fieldName, r.filter, opts...)

	if startR != nil {
		logger.GetLogger().AddTraceInfoContextRequest(ctx).Info().Dur("latency", time.Since(*startR)).Msg("mongodb_latency: DistinctDocs.Distinct")
	}

	return rs, err
}

func (r *Repository[T]) SetLimit(limit int64) *Repository[T] {
	r.optsFind.Limit = &limit
	return r
}

func (r *Repository[T]) SetSkip(skip int64) *Repository[T] {
	r.optsFind.Skip = &skip
	return r
}

func (r *Repository[T]) SetSkipOne(skip int64) *Repository[T] {
	r.optsFindOne.Skip = &skip
	return r
}

func (r *Repository[T]) SetProjection(projection bson.M) *Repository[T] {
	r.optsFind.Projection = projection
	return r
}

func (r *Repository[T]) SetProjectionOne(projection bson.M) *Repository[T] {
	r.optsFindOne.Projection = projection
	return r
}

func (r *Repository[T]) SetHint(hint bson.M) *Repository[T] {
	r.optsFind.Hint = hint
	return r
}

func (r *Repository[T]) SetHintOne(hint bson.M) *Repository[T] {
	r.optsFindOne.Hint = hint
	return r
}

func (r *Repository[T]) SetMetricMethod(method string) *Repository[T] {
	r.metricMethod = method
	return r
}

func (r *Repository[T]) filterEncrypt() error {
	if r.keyEncrypt == "" || len(r.fieldsNameEnc) == 0 {
		return nil
	}

	for i, fil := range r.filter {
		if _, ok := r.fieldsNameEnc[fil.Key]; ok {
			if data, _ok := fil.Value.(string); _ok {
				enc, err := utils.Encrypt(data, r.keyEncrypt)
				if err != nil {
					return fmt.Errorf("filter encrypt error: %v", err)
				}
				r.filter[i].Value = enc
			}
		}
	}

	return nil
}
