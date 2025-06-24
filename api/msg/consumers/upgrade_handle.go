// package tier

// import (
// 	"context"
// 	"encoding/json"
// )

// type CalculateTier struct {
// 	db          *mongodb.DatabaseStorage
// 	log         *logger.Logger
// 	tierService services.ITierService
// }

// func NewCalculateTier(tierService services.ITierService) *CalculateTier {
// 	return &CalculateTier{
// 		log:         logger.GetLogger(),
// 		tierService: tierService,
// 	}
// }

// func (s *CalculateTier) Handle(ctx context.Context, data interface{}) error {
// 	ev := models.MessageTier{}
// 	byteData, err := json.Marshal(data)
// 	if err != nil {
// 		s.log.Err(err).Msg("decode sync data tier upgrade message failed")
// 		return err
// 	}
// 	err = json.Unmarshal(byteData, &ev)
// 	if err != nil {
// 		s.log.Err(err).Msg("decode sync data tier upgrade failed")
// 		return err
// 	}

// 	dataService := adapters.AdapterTier{}.ConvertRequest2ServiceTransaction(ev)
// 	err = s.tierService.SyncDataTier(ctx, dataService)
// 	if err != nil {
// 		s.log.Err(err).Msg("excuse logic sync data tier upgrade failed")
// 		return err
// 	}

// 	return nil
// }
