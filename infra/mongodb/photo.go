package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"io"
	"log"
	"time"
)

type Photo struct {
	Title       string    `bson:"title"`
	ContentType string    `bson:"content_type"`
	Size        int64     `bson:"size"`
	UploadedAt  time.Time `bson:"uploaded_at"`
	UserID      uint64    `bson:"user_id"`
}

func UploadPhoto(userID uint64, title, contentType string, r io.Reader, size int64) (string, error) {
	if Bucket == nil {
		log.Fatalf("Bucket is not initialized.")
		return "", nil
	}

	photo := Photo{
		Title:       title,
		ContentType: contentType,
		Size:        size,
		UploadedAt:  time.Now(),
		UserID:      userID,
	}
	uploadOpts := options.GridFSUpload().SetMetadata(photo)

	objectID, err := Bucket.UploadFromStream(
		context.TODO(),
		photo.Title,
		r,
		uploadOpts,
	)
	if err != nil {
		log.Fatalf("Error while uploading photo: %v", err)
		return "", err
	}

	return objectID.Hex(), nil
}
