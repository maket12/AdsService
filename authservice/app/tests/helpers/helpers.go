package helpers

import (
	"ads/authservice/domain/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func MakeAccessClaims() *entity.AccessClaims {
	return &entity.AccessClaims{
		Type:   "access",
		UserID: 1,
		Email:  "",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(1),
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func MakeRefreshClaims(jti string, exp time.Time, iat time.Time) *entity.RefreshClaims {
	return &entity.RefreshClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(iat),
		},
	}
}

func MakeTestUser(email, password, role string) *entity.User {
	passwordHashes := map[string]string{
		"password123": "$2a$10$SH7gGMiqW3s.spBmBh1ZgeK4QHfenm8K/W/E680uGCXTd1.kQP1OW",
		"password178": "$2a$10$K.ssBhqJ9wEyBKjPqhE6z.G/Oqd9yJrZMQd5EUhLgVduCz7cXiHjO",
		"newpass":     "$2a$10$LoLYO8YRBeb6XXH96oeD2OWkj5ywIQgYiZFs9ZCs95hITSHfbAXpa",
	}
	return &entity.User{
		ID:       1,
		Email:    email,
		Password: passwordHashes[password],
		Role:     role,
	}
}
