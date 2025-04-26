package adapters

import (
	"go-source/internal/domains"
	"go-source/repositories/entity"
)

type AdapterEntity struct {
}

func (a AdapterEntity) ConvertServToRepo(data *entity.Entity) *domains.Entity {
	return &domains.Entity{
		ID:     data.Id.Hex(),
		Status: data.Status,
	}
}
