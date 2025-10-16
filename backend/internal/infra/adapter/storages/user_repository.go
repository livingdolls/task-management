package storages

import (
	"errors"
	"task-management/internal/applications/ports/repository"
	"task-management/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create implements repository.UserRepository.
func (u *userRepository) Create(user *domain.User) error {
	return u.db.Create(user).Error
}

// FindByID implements repository.UserRepository.
func (u *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

// FindByUsername implements repository.UserRepository.
func (u *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
