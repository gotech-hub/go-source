package config

import (
	"github.com/caarlos0/env/v7"
)

type SystemConfig struct {
	Env      string `env:"ENV,required,notEmpty"`
	HttpPort uint64 `env:"HTTP_PORT,required,notEmpty"`
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
