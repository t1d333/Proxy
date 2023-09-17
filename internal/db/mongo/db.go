package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitDB(ctx context.Context, logger *zap.SugaredLogger) *mongo.Client {
	logger.Info("trying connect to mongo...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		logger.Fatal("mongo client creation failed", zap.Error(err))
		return nil
	}

	logger.Info("trying ping mongo...")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal("mongo is not available", zap.Error(err))
	}

	logger.Info("connection to mongo successfully")
	return client
}
