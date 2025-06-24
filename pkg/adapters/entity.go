package adapters

import (
	"go-source/internal/domains"
	entity "go-source/repositories/entity1"
)

type AdapterEntity struct {
}

func (a AdapterEntity) ConvertServToRepo(data *entity.Entity) *domains.Entity {
	return &domains.Entity{
		ID:     data.Id.Hex(),
		Status: data.Status,
	}
}
