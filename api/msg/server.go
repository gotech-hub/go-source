// package msg

// import (
// 	"context"
// 	"go-source/config"
// 	logger "go-source/pkg/log"
// 	"go-source/pkg/queue/kafka"
// )

// type msgBroker struct {
// 	conf *config.SystemConfig
// }

// func NewMsgBroker(conf *config.SystemConfig) *msgBroker {
// 	return &msgBroker{
// 		conf: conf,
// 	}
// }

// func (app *msgBroker) Start(ctx context.Context) {
// 	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)
// 	// Consumer for calculate user tier upgrade
// 	csTier := kafka.NewConsumer(app.conf.KafkaConfig, []string{app.conf.KafkaTopicConfig.EarnPointSuccess})
// 	csTierSv := tier.NewConsumerTier(csTier, app.tierService)
// 	go func() {
// 		err := csTierSv.Start(ctx)
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("csTierUpgrade failed")
// 		}
// 	}()
// }
