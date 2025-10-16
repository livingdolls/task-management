package services

import "task-management/internal/domain"

type AuthService interface {
	Register(name, username, password string) (*domain.User, error)
	Login(username, password string) (string, *domain.User, error)
	Me(userID uint) (*domain.User, error)
}
