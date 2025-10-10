package profilerepo

import (
	"ads/internal/adapters/out/postgres/shared"
	"ads/internal/core/domain/model/profile"
	"ads/internal/pkg/errs"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type ProfilesRepo struct {
	tracker shared.Tracker
}

func NewProfilesRepo(tracker shared.Tracker) (*ProfilesRepo, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	return &ProfilesRepo{tracker: tracker}, nil
}

func (r *ProfilesRepo) AddProfile(ctx context.Context, aggregate *profile.Profile) (*profile.Profile, error) {
	r.tracker.Track(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(aggregate).Error
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created new profile: %v", aggregate.UserID())

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
