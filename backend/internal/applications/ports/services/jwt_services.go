package services

import "task-management/internal/domain"

type JWTService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (*domain.JWTClaims, error)
}
