package mongodb

import (
	"AdsService/userservice/config"
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log/slog"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
	Bucket *mongo.GridFSBucket
)

func InitMongoDB(cfg *config.Config, logger *slog.Logger) error {
	logger.Info("connecting to MongoDB...",
		slog.String("uri", cfg.MongoURI),
		slog.String("database", cfg.MongoDB),
	)

	var err error
	Client, err = mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return err
	}

	DB = Client.Database(cfg.MongoDB)
	Bucket = DB.GridFSBucket(options.GridFSBucket().SetName(cfg.MongoBucket))

	logger.Info("✅ MongoDB connected successfully")
	return nil
}

func CloseMongoDB(logger *slog.Logger) {
	if Client != nil {
		_ = Client.Disconnect(context.Background())
		logger.Info("MongoDB connection closed")
	}
}
