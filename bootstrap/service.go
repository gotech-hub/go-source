package bootstrap

import "go-source/pkg/database/redis"

type Services struct {
}

func NewServices(repo *Repositories, redis *redis.Client, clients *Clients) *Services {
	service := &Services{}
	return service
}
