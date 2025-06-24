package bootstrap

import (
	"context"
	"go-source/config"
	"go-source/pkg/database/mongodb"
	logger "go-source/pkg/log"
)

var (
	databaseConnection *DatabaseConnection
)

type DatabaseConnection struct {
	Connection *mongodb.DatabaseStorage
}

func NewDatabaseConnection(ctx context.Context) *DatabaseConnection {
	log := logger.GetLogger().AddTraceInfoContextRequest(ctx)

	connection, err := mongodb.ConnectMongoDB(ctx, &config.GetInstance().MongoDBConfig)
	handleError(log, err, "Connect MongoDB promotion failed!")

	databaseConnection = &DatabaseConnection{
		Connection: connection,
	}

	return databaseConnection
}

func handleError(log *logger.Logger, err error, errMsg string) {
	if err != nil {
		log.Fatal().Msgf(errMsg, err)
	}
}
