package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FilterPlayer struct {
	filter bson.D

	optsFind options.FindOptions
	sort     bson.D

	optsFindOne options.FindOneOptions
	sortOne     bson.D

	metricComponent string
	metricMethod    string
}

type FilterPlayerOption func(*FilterPlayer)

func WithMetricComponent(component string) FilterPlayerOption {
	return func(f *FilterPlayer) {
		f.metricComponent = component
	}
}

func WithMetricMethod(method string) FilterPlayerOption {
	return func(f *FilterPlayer) {
		f.metricMethod = method
	}
}

func defaultFilterPlayer() *FilterPlayer {
	return &FilterPlayer{
		filter:      bson.D{},
		optsFind:    options.FindOptions{},
		sort:        bson.D{},
		optsFindOne: options.FindOneOptions{},
		sortOne:     bson.D{},
	}
}

func NewFilterPlayer(opts ...FilterPlayerOption) *FilterPlayer {
	filterPlayer := defaultFilterPlayer()

	for _, opt := range opts {
		opt(filterPlayer)
	}

	return filterPlayer
}

// Append Not support bson.A
func (f *FilterPlayer) Append(data interface{}) *FilterPlayer {
	// check type data
	switch data.(type) {
	case bson.D:
		f.filter = append(f.filter, data.(bson.D)...)
	case bson.E:
		s := data.(bson.E)
		f.filter = append(f.filter, bson.E{Key: s.Key, Value: s.Value})
	case bson.M:
		for k, v := range data.(bson.M) {
			f.filter = append(f.filter, bson.E{Key: k, Value: v})
		}
	}
	return f
}

// AppendSort Not support bson.A
func (f *FilterPlayer) AppendSort(sort interface{}) *FilterPlayer {
	// check type sort
	switch sort.(type) {
	case bson.D:
		f.sort = append(f.sort, sort.(bson.D)...)
	case bson.E:
		s := sort.(bson.E)
		f.sort = append(f.sort, bson.E{Key: s.Key, Value: s.Value})
	case bson.M:
		for k, v := range sort.(bson.M) {
			f.sort = append(f.sort, bson.E{Key: k, Value: v})
		}
	}
	return f
}

// AppendSortOne Not support bson.A
func (f *FilterPlayer) AppendSortOne(sort interface{}) *FilterPlayer {
	// check type sort
	switch sort.(type) {
	case bson.D:
		f.sortOne = append(f.sortOne, sort.(bson.D)...)
	case bson.E:
		s := sort.(bson.E)
		f.sortOne = append(f.sortOne, bson.E{Key: s.Key, Value: s.Value})
	case bson.M:
		for k, v := range sort.(bson.M) {
			f.sortOne = append(f.sortOne, bson.E{Key: k, Value: v})
		}
	}
	return f
}

func (f *FilterPlayer) SetMethod(method string) *FilterPlayer {
	f.metricMethod = method
	return f
}
