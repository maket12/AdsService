package mongo

import (
	"AdsService/userservice/domain/entity"
	"AdsService/userservice/domain/port"
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"io"
	"log"
	"time"
)

type PhotoRepo struct {
	Bucket *mongo.GridFSBucket
}

func NewPhotoRepo(bucket *mongo.GridFSBucket) port.PhotoRepository {
	return &PhotoRepo{Bucket: bucket}
}

func (r *PhotoRepo) UploadPhoto(userID uint64, title, contentType string, rdr io.Reader, size int64) (string, error) {
	if r.Bucket == nil {
		log.Fatalf("Bucket is not initialized.")
		return "", nil
	}

	photo := entity.Photo{
		Title:       title,
		ContentType: contentType,
		Size:        size,
		UploadedAt:  time.Now(),
		UserID:      userID,
	}
	uploadOpts := options.GridFSUpload().SetMetadata(photo)

	objectID, err := r.Bucket.UploadFromStream(
		context.TODO(),
		photo.Title,
		rdr,
		uploadOpts,
	)
	if err != nil {
		log.Fatalf("Error while uploading photo: %v", err)
		return "", err
	}

	return objectID.Hex(), nil
}
