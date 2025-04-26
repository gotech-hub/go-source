package kafka

type KafkaConfig struct {
	BootstrapServers string `env:"BOOTSTRAP_SERVERS"`
	GroupID          string `env:"GROUP_ID"`
	AutoOffsetReset  string `env:"AUTO_OFFSET_RESET"`
	SecurityProtocol string `env:"SECURITY_PROTOCOL"`
	SaslMechanism    string `env:"SASL_MECHANISM"`
	SaslMechanisms   string `env:"SASL_MECHANISMS"`
	SaslUsername     string `env:"SASL_USERNAME"`
	SaslPassword     string `env:"SASL_PASSWORD"`
}
