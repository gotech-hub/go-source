package bootstrap

import "go-source/pkg/database/mongodb"

var (
	repositories *Repositories
)

type Repositories struct {
	// profile

}

func NewRepositories(db *mongodb.DatabaseStorage) *Repositories {
	repositories = &Repositories{}

	return repositories
}
