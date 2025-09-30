package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
	Bucket *mongo.GridFSBucket
)

func InitMongoDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	var err error
	Client, err = mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGODB_DB_NAME is not set")
	}

	bucketName := os.Getenv("MONGODB_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("MONGODB_BUCKET_NAME is not set")
	}

	DB = Client.Database(dbName)
	Bucket = DB.GridFSBucket(options.GridFSBucket().SetName(bucketName))

	log.Println("MongoDB connected successfully")
}

func CloseMongoDB() {
	if Client != nil {
		_ = Client.Disconnect(context.Background())
		log.Println("MongoDB connection closed")
	}
}
