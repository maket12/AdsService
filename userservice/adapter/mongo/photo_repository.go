package mongo

import (
	"AdsService/userservice/domain/entity"
	"AdsService/userservice/domain/port"
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"io"
	"log/slog"
	"time"
)

type PhotoRepo struct {
	Bucket *mongo.GridFSBucket
	logger *slog.Logger
}

func NewPhotoRepo(bucket *mongo.GridFSBucket, logger *slog.Logger) port.PhotoRepository {
	return &PhotoRepo{
		Bucket: bucket,
		logger: logger,
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
	r.logger.Info("Successfully uploaded photo of user[%v] with hexID[%v]", userID, hexID)

	return hexID, nil
}
