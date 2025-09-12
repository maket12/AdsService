package graph

import (
	"AdsService/infra/gateway/graph/model"
	userpb "AdsService/userservice/proto"
	"fmt"
	"time"
)

func toTimePtr(t time.Time) *time.Time {
	return &t
}

func mapUserProfile(p *userpb.Profile) *model.UserProfile {
	if p == nil {
		return nil
	}
	return &model.UserProfile{
		UserID:               fmt.Sprint(p.UserId),
		Name:                 &p.Name,
		Phone:                &p.Phone,
		PhotoID:              &p.PhotoId,
		NotificationsEnabled: p.NotificationsEnabled,
		Subscriptions:        p.Subscriptions,
		UpdatedAt:            toTimePtr(p.UpdatedAt.AsTime()),
	}
}

func getOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
