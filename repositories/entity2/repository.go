package entity

import (
	"context"
	"errors"
	"go-source/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type EntityRepository struct {
	*mongodb.Repository[Entity]
}

type IEntityRepository interface {
	Create(ctx context.Context, data *Entity) error
	Get(ctx context.Context, id primitive.ObjectID) (*Entity, error)
}

var (
	instance *EntityRepository
	once     sync.Once
)

func NewEntityRepository(dbStorage *mongodb.DatabaseStorage) IEntityRepository {
	once.Do(func() {
		instance = &EntityRepository{
			Repository: mongodb.NewRepository[Entity](dbStorage),
		}
	})
	return instance
}

func (r *EntityRepository) R() *EntityRepository {
	return &EntityRepository{
		Repository: r.Repository.NewFilterPlayer(),
	}
}

func (r *EntityRepository) ById(id primitive.ObjectID) *EntityRepository {
	filter := bson.M{
		FEntityId: id,
	}
	r.Append(filter)
	return r
}

func (r *EntityRepository) Create(ctx context.Context, data *Entity) error {
	_, err := r.R().CreateOneDocument(ctx, data)
	if !errors.Is(err, nil) {
		return err
	}
	return nil
}

func (r *EntityRepository) Get(ctx context.Context, id primitive.ObjectID) (*Entity, error) {
	rs, err := r.R().ById(id).FindOneDoc(ctx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
