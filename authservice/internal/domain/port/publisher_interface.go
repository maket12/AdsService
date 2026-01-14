package port

import (
	"ads/authservice/internal/domain/model"
	"context"
)

type AccountPublisher interface {
	PublishAccountCreate(ctx context.Context, event *model.AccountCreatedEvent) error
}
