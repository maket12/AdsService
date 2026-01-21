package graph

import (
	"ads/pkg/generated/auth_v1"
	"ads/pkg/generated/user_v1"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthClient auth_v1.AuthServiceClient
	UserClient user_v1.UserServiceClient
}
