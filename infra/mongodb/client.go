package mongodb

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
	Bucket      *mongo.GridFSBucket
	BucketName  string
	DbName      string
)

func InitMongoDB() {
	err1 := godotenv.Load()
	if err1 != nil {
		return
	}

	uri := os.Getenv("MONGODB_URI")
	var err error
	MongoClient, err = mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error while connecting to the MongoDB: %v", err)
	}

	if DbName = os.Getenv("MONGODB_DB_NAME"); DbName == "" {
		log.Fatalf("Database name is not found")
	}
	if BucketName = os.Getenv("MONGODB_BUCKET_NAME"); BucketName == "" {
		log.Fatalf("Collection name is not found")
	}

	MongoDB = MongoClient.Database(DbName)
	bucketOpts := options.GridFSBucket().SetName(BucketName)

	Bucket = MongoDB.GridFSBucket(bucketOpts)
}

func CloseMongoDB() {
	if MongoClient != nil {
		_ = MongoClient.Disconnect(context.Background())
	}
}
