package mongo

import (
	"ads/userservice/domain/entity"
	"ads/userservice/domain/port"
	"context"
	"io"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PhotoRepo struct {
	Bucket *mongo.GridFSBucket
	logger *slog.Logger
}

func NewPhotoRepo(bucket *mongo.GridFSBucket, log *slog.Logger) port.PhotoRepository {
	return &PhotoRepo{
		Bucket: bucket,
		logger: log,
	}
}

func (r *PhotoRepo) UploadPhoto(ctx context.Context, userID uint64, title, contentType string, rdr io.Reader, size int64) (string, error) {
	photo := entity.Photo{
		Title:       title,
		ContentType: contentType,
		Size:        size,
		UploadedAt:  time.Now(),
		UserID:      userID,
	}
	uploadOpts := options.GridFSUpload().SetMetadata(photo)

	objectID, err := r.Bucket.UploadFromStream(
		ctx,
		photo.Title,
		rdr,
		uploadOpts,
	)
	if err != nil {
		return "", err
	}

	hexID := objectID.Hex()
	r.logger.InfoContext(ctx, "Successfully uploaded photo of user[%v] with hexID[%v]", userID, hexID)

	return hexID, nil
}
