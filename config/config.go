package config

import (
	"go-source/pkg/database/mongodb"
	"go-source/pkg/database/redis"

	"github.com/caarlos0/env/v7"
)

type SystemConfig struct {
	Env            string `env:"ENV,required,notEmpty"`
	HttpPort       uint64 `env:"HTTP_PORT,required,notEmpty"`
	ServiceName    string `env:"SERVICE_NAME,required,notEmpty"`
	ServiceVersion string `env:"SERVICE_VERSION,required,notEmpty"`

	MongoDBConfig mongodb.MongoDBConfig `envPrefix:"MONGO_DB_" envSeparator:"_"`
	RedisConfig   redis.RedisConfig     `envPrefix:"REDIS_" envSeparator:"_"`
}

var configSingletonObj *SystemConfig

func LoadConfig() (cf *SystemConfig, err error) {

	if configSingletonObj != nil {
		cf = configSingletonObj
		return
	}

	cf = &SystemConfig{}
	if err = env.Parse(cf); err != nil {
		return
	}

	configSingletonObj = cf
	return
}

func GetInstance() *SystemConfig {
	return configSingletonObj
}
