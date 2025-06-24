// package tier

// import (
// 	"context"
// 	"encoding/json"
// 	logger "gitlab.gt.vng.vn/vng-rewards/share/vgr-be-common/log"
// 	"gitlab.gt.vng.vn/vng-rewards/share/vgr-be-common/queue/kafka"
// 	"vgr-tier/api/msg/models"
// 	"vgr-tier/internal/services"
// 	constant "vgr-tier/pkg/helpers/constants"
// )

// type ConsumerTier struct {
// 	log                 *logger.Logger
// 	cs                  *kafka.Consumer
// 	handleCalculateTier *CalculateTier
// }

// func NewConsumerTier(
// 	cs *kafka.Consumer,
// 	tierService services.ITierService,
// ) *ConsumerTier {

// 	return &ConsumerTier{
// 		cs:                  cs,
// 		log:                 logger.GetLogger(),
// 		handleCalculateTier: NewCalculateTier(tierService),
// 	}
// }

// func (s *ConsumerTier) Start(ctx context.Context) error {
// 	s.cs.OnEvent(func(ctx context.Context, key, value []byte) error {
// 		var ev models.ConsumerMessage
// 		if err := json.Unmarshal(value, &ev); err != nil {
// 			s.log.Err(err).Str("key", string(key)).Str("value", string(value)).Msg("decode failed")
// 			return nil // ignore this event
// 		}
// 		s.log.Info().Any("Calculate user tier", ev).Msg("event data")

// 		switch ev.EventType {
// 		case constant.TierUpgrade, constant.TierRefund:
// 			err := s.handleCalculateTier.Handle(ctx, ev.Data)
// 			if err != nil {
// 				s.log.Err(err).Msg("handle calculate user tier failed")
// 			}
// 		}

// 		return nil
// 	})

// 	if err := s.cs.Start(ctx); err != nil {
// 		return err
// 	}

// 	return nil
// }
