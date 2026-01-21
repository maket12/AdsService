package port

import (
	"ads/authservice/internal/domain/entity"
	"context"
)

type AccountPublisher interface {
	PublishAccountCreate(ctx context.Context, event *entity.AccountCreatedEvent) error
}
