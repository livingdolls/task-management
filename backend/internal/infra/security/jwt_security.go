package security

import (
	"task-management/internal/applications/ports/services"
	"task-management/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAdapter struct {
	secret string
	ttl    time.Duration
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTAdapter(secret string, ttl time.Duration) services.JWTService {
	return &JWTAdapter{
		secret: secret,
		ttl:    ttl,
	}
}

// GenerateToken implements services.JWTService.
func (j *JWTAdapter) GenerateToken(user *domain.User) (string, error) {
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "task-management-services",
			Subject:   "user-authentication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken implements services.JWTService.
func (j *JWTAdapter) ValidateToken(token string) (*domain.User, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &domain.User{
		ID:       claims.UserID,
		Username: claims.Username,
	}, nil
}
