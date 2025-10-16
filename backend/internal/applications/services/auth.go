package services

import (
	"errors"
	"task-management/internal/applications/ports/repository"
	"task-management/internal/applications/ports/services"
	"task-management/internal/domain"
	"task-management/internal/utils"
)

var (
	ErrUserExists      = errors.New("username is already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

type authService struct {
	repo repository.UserRepository
	jwt  services.JWTService
}

func NewAuthService(repo repository.UserRepository, jwt services.JWTService) services.AuthService {
	return &authService{
		repo: repo,
		jwt:  jwt,
	}
}

// Login implements services.AuthService.
func (a *authService) Login(username string, password string) (string, *domain.User, error) {
	user, err := a.repo.FindByUsername(username)
	if err != nil {
		return "", nil, err
	}

	if user == nil || utils.CheckPassword(user.Password, password) != nil {
		return "", nil, ErrInvalidPassword
	}

	token, err := a.jwt.GenerateToken(user)

	if err != nil {
		return "", nil, err
	}

	user.Password = ""
	return token, user, nil
}

// Register implements services.AuthService.
func (a *authService) Register(name string, username string, password string) (*domain.User, error) {
	// check if username already exists
	existingUser, err := a.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     name,
		Username: username,
		Password: hashedPassword,
	}

	if err := a.repo.Create(user); err != nil {
		return nil, err
	}

	// hide password, jangan kirim ke handler
	user.Password = ""

	return user, nil
}

func (a *authService) Me(userID uint) (*domain.User, error) {
	user, err := a.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserExists
	}

	user.Password = ""
	return user, nil
}
