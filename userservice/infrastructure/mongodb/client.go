package mongodb

import (
	"ads/userservice/pkg"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoData struct {
	Client *mongo.Client
	DB     *mongo.Database
	Bucket *mongo.GridFSBucket
}

func InitMongoDB(cfg *pkg.Config) (*MongoData, error) {
	var (
		mongoData MongoData
		err       error
	)

	fmt.Printf("connecting to MongoDB...  uri=%s database_name=%s\n",
		cfg.MongoURI, cfg.MongoDB)

	mongoData.Client, err = mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	mongoData.DB = mongoData.Client.Database(cfg.MongoDB)
	mongoData.Bucket = mongoData.DB.GridFSBucket(options.GridFSBucket().SetName(cfg.MongoBucket))

	fmt.Print("✅ MongoDB connected successfully")
	return &mongoData, nil
}

func CloseMongoDB(mongoData *MongoData) {
	if mongoData.Client != nil {
		_ = mongoData.Client.Disconnect(context.Background())
		fmt.Print("MongoDB connection closed")
	}
}
