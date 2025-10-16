package storages

import "gorm.io/gorm"

type taskRepository struct {
	db *gorm.DB
}
