package entity

import "time"

type Photo struct {
	Title       string    `bson:"title"`
	ContentType string    `bson:"content_type"`
	Size        int64     `bson:"size"`
	UploadedAt  time.Time `bson:"uploaded_at"`
	UserID      uint64    `bson:"user_id"`
}
