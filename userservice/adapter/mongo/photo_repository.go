package mongo

import (
	"ads/userservice/domain/entity"
	"ads/userservice/domain/port"
	"context"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PhotoRepo struct {
	Bucket *mongo.GridFSBucket
}

func NewPhotoRepo(bucket *mongo.GridFSBucket) port.PhotoRepository {
	return &PhotoRepo{
		Bucket: bucket,
	}
}

func (r *PhotoRepo) UploadPhoto(userID uint64, title, contentType string, rdr io.Reader, size int64) (string, error) {
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
		return "", err
	}

	hexID := objectID.Hex()
	fmt.Printf("Successfully uploaded photo of user[%v] with hexID[%v]", userID, hexID)

	return hexID, nil
}
