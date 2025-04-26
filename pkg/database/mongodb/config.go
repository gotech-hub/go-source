package mongodb

type MongoDBConfig struct {
	DatabaseURI          string `env:"DATABASE_URI,required,notEmpty"`
	DatabaseName         string `env:"DATABASE_NAME,required,notEmpty"`
	IsEnableDebugLogger  bool   `env:"IS_ENABLE_DEBUG_LOGGER"`
	ShouldMeasureLatency bool   `env:"SHOULD_MEASURE_LATENCY"`
}

type MultiConnMongoConfig map[string]map[string]string
