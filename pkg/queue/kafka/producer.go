package kafka

import (
	"context"
	"encoding/json"
	logger "go-source/pkg/log"
	"go-source/pkg/utils"
	"hash/crc32"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerInterface interface {
	Publish(ctx context.Context, key, value interface{}) error
}

type Producer struct {
	pr            *kafka.Producer
	topic         string
	numPartitions int32
}

func NewProducer(cfg KafkaConfig, topic string, numPartitions ...int32) *Producer {
	log := logger.GetLogger()

	cfgMap := kafka.ConfigMap{
		"bootstrap.servers":  cfg.BootstrapServers,
		"enable.idempotence": true,
		"acks":               "all",
	}

	if cfg.SecurityProtocol != "" {
		cfgMap["security.protocol"] = cfg.SecurityProtocol
	}

	if cfg.SaslMechanism != "" {
		cfgMap["sasl.mechanism"] = cfg.SaslMechanism
	} else if cfg.SaslMechanisms != "" {
		cfgMap["sasl.mechanisms"] = cfg.SaslMechanisms
	}

	if cfg.SaslUsername != "" {
		cfgMap["sasl.username"] = cfg.SaslUsername
	}

	if cfg.SaslPassword != "" {
		cfgMap["sasl.password"] = cfg.SaslPassword
	}

	pr, err := kafka.NewProducer(&cfgMap)

	if err != nil {
		log.Fatal().Err(err).Msg("init kafka producer failed")
	}

	log.Info().Msgf("init kafka producer success : TOPIC = %v", topic)

	res := &Producer{
		pr:    pr,
		topic: topic,
	}

	if len(numPartitions) > 0 {
		res.numPartitions = numPartitions[0]
	}

	return res
}

func marshal(val interface{}) ([]byte, error) {
	switch v := val.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(val)
	}
}

func (s *Producer) Publish(ctx context.Context, key, value interface{}) error {
	keyData, err := marshal(key)
	if err != nil {
		return err
	}
	valueData, err := marshal(value)
	if err != nil {
		return err
	}

	var header []byte
	traceInfo := utils.GetRequestIdByContext(ctx)
	if traceInfo != nil {
		header, err = marshal(traceInfo)
		if err != nil {
			return err
		}
	}

	msg := &kafka.Message{
		Key:   keyData,
		Value: valueData,
		Headers: []kafka.Header{
			{
				Key:   utils.KeyTraceInfo,
				Value: header,
			},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: kafka.PartitionAny,
		},
	}
	if err = s.pr.Produce(msg, nil); err != nil {
		return err
	}
	return nil
}

func (s *Producer) PublishWithPartition(ctx context.Context, key, value interface{}, partition int32) error {
	keyData, err := marshal(key)
	if err != nil {
		return err
	}
	valueData, err := marshal(value)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		Key:   keyData,
		Value: valueData,
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: partition,
		},
	}
	if err = s.pr.Produce(msg, nil); err != nil {
		return err
	}
	return nil
}

// PublishWithPartitionCRC32 publish message with partition is crc32(key) % numPartitions
func (s *Producer) PublishWithPartitionCRC32(ctx context.Context, key, value interface{}) error {
	keyData, err := marshal(key)
	if err != nil {
		return err
	}

	valueData, err := marshal(value)
	if err != nil {
		return err
	}

	var header []byte
	traceInfo := utils.GetRequestIdByContext(ctx)
	if traceInfo != nil {
		header, err = marshal(traceInfo)
		if err != nil {
			return err
		}
	}

	partition := kafka.PartitionAny
	if s.numPartitions > 0 {
		partition = int32(crc32.ChecksumIEEE(keyData)) % s.numPartitions
		if partition < 0 {
			partition = -partition
		}
	}

	msg := &kafka.Message{
		Key:   keyData,
		Value: valueData,
		Headers: []kafka.Header{
			{
				Key:   utils.KeyTraceInfo,
				Value: header,
			},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: partition,
		},
	}
	if err = s.pr.Produce(msg, nil); err != nil {
		return err
	}

	return nil
}

func (s *Producer) PublishBytes(ctx context.Context, key, value []byte) error {
	msg := &kafka.Message{
		Key:   key,
		Value: value,
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topic,
			Partition: kafka.PartitionAny,
		},
	}
	return s.pr.Produce(msg, nil)
}

func (s *Producer) PublishMessage(ctx context.Context, msg *kafka.Message) error {
	return s.pr.Produce(msg, nil)
}

func (s *Producer) PublishWithTopic(ctx context.Context, topic string, key, value interface{}) error {
	keyData, err := marshal(key)
	if err != nil {
		return err
	}
	valueData, err := marshal(value)
	if err != nil {
		return err
	}

	var header []byte
	traceInfo := utils.GetRequestIdByContext(ctx)
	if traceInfo != nil {
		header, err = marshal(traceInfo)
		if err != nil {
			return err
		}
	}

	msg := &kafka.Message{
		Value: valueData,
		Key:   keyData,
		Headers: []kafka.Header{
			{
				Key:   utils.KeyTraceInfo,
				Value: header,
			},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}

	if err := s.pr.Produce(msg, nil); err != nil {
		return err
	}
	return nil
}

func (s *Producer) GetTopicName() string {
	return s.topic
}
