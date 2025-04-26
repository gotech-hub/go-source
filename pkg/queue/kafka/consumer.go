package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	logger "go-source/pkg/log"
	"go-source/pkg/utils"
	"sync"
	"time"
)

var (
	ErrAlreadyStarted  = errors.New("already started")
	ErrNilEventHandler = errors.New("event handlers is nil")
)

type OnEventHandler func(ctx context.Context, key, value []byte) error

type ConsumerInterface interface {
	OnEvent(handler OnEventHandler)
	Start(ctx context.Context) error
	Shutdown(ctx context.Context)
}

type Consumer struct {
	cs      *kafka.Consumer
	topics  []string
	handler OnEventHandler

	started bool
	mu      sync.RWMutex
}

func NewConsumer(cfg KafkaConfig, topics []string) *Consumer {
	log := logger.GetLogger()

	cfgMap := kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": cfg.AutoOffsetReset,
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

	cs, err := kafka.NewConsumer(&cfgMap)

	if err != nil {
		log.Fatal().Err(err).Msg("init kafka consumer failed")
	}

	log.Info().Msgf("init kafka consumer success : TOPIC = %v", topics)
	return &Consumer{
		cs:     cs,
		topics: topics,
	}
}

func (s *Consumer) OnEvent(handler OnEventHandler) {
	s.mu.Lock()
	if handler != nil {
		s.handler = handler
	}
	s.mu.Unlock()
}

func (s *Consumer) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return ErrAlreadyStarted
	}
	if s.handler == nil {
		s.mu.Unlock()
		return ErrNilEventHandler
	}
	s.started = true
	s.mu.Unlock()

	if err := s.cs.SubscribeTopics(s.topics, nil); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		log := logger.GetLogger()
		msg, err := s.cs.ReadMessage(10 * time.Second)
		if err != nil {
			var kerr kafka.Error
			if errors.As(err, &kerr); kerr.Code() == kafka.ErrTimedOut {
				continue
			}
			log.Warn().Err(err).Msg("kafka read message failed")
			continue
		}

		traceInfoExisted := false
		newCtx := context.Background()

		for _, h := range msg.Headers {
			if h.Key == utils.KeyTraceInfo {
				traceInfo := utils.TraceInfo{}
				if err = json.Unmarshal(h.Value, &traceInfo); err != nil {
					break
				}
				newCtx = context.WithValue(newCtx, utils.KeyTraceInfo, traceInfo)
				traceInfoExisted = true
			}
		}

		if !traceInfoExisted {
			newCtx, _ = utils.NewContextWithRequestId(ctx)
		}

		log = log.AddTraceInfoContextRequest(newCtx)
		log.Info().
			Interface("topicPartition", msg.TopicPartition).
			Str("value", string(msg.Value)).
			Str("key", string(msg.Key)).
			Time("timestamp", msg.Timestamp).
			Int("timestampType", int(msg.TimestampType)).
			Interface("opaque", msg.Opaque).
			Interface("headers", msg.Headers).
			Msg("kafka read message success")

		if err = s.handler(newCtx, msg.Key, msg.Value); err != nil {
			log.Err(err).Msg("kafka handlers failed")
			continue
		}

		if _, err = s.cs.CommitMessage(msg); err != nil {
			log.Err(err).Msg("kafka commit failed")
		}
	}
}

func (s *Consumer) Shutdown(ctx context.Context) {
}

func (s *Consumer) GetTopics() []string {
	return s.topics
}
