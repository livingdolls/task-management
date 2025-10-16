package security

import (
	"errors"
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
func (j *JWTAdapter) ValidateToken(token string) (*domain.JWTClaims, error) {
	t, err := jwt.ParseWithClaims(token, &domain.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("invalid token")
	}

	claims, ok := t.Claims.(*domain.JWTClaims)

	if !ok || !t.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
