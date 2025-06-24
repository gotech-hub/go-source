package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Entity struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Status    string             `bson:"status"` // pending, accepted, blocked
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (Entity) CollectionName() string {
	return ColEntity
}

func (Entity) IndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{}
}
