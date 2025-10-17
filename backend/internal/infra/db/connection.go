package db

import (
	"database/sql"
	"fmt"
	"task-management/internal/config"
	"task-management/internal/domain"
	"task-management/internal/infra/logger"
	"task-management/internal/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	gormConfig := &gorm.Config{}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	// jalankan migrasi otomatis
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Task{},
	)

	if err != nil {
		logger.Error("Failed to run auto migration", zap.Error(err))
		return nil, fmt.Errorf("failed to run auto migration: %w", err)
	}

	// buat user default jika belum ada
	createDefaultUser(db)

	sqlDB.SetMaxOpenConns(cfg.MaxOpenCons)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleCons)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Minute)

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.String("name", cfg.Name),
	)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	if d.DB != nil {
		return nil
	}

	sqlDB, err := d.DB.DB()

	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil && err != sql.ErrConnDone {
		return fmt.Errorf("failed to close DB connection: %w", err)
	}

	logger.Info("Database connection closed")

	return nil
}

func createDefaultUser(db *gorm.DB) {
	const defaultUsername = "admin"
	const defaultPassword = "admin123"

	var user domain.User
	result := db.First(&user, "username = ?", defaultUsername)

	if result.Error == gorm.ErrRecordNotFound {
		hashed, _ := utils.HashPassword(defaultPassword)

		newUser := domain.User{
			Username: defaultUsername,
			Password: string(hashed),
		}

		if err := db.Create(&newUser).Error; err != nil {
			logger.Info("failed to create default user")
		} else {
			logger.Info("Default user created",
				zap.String("username", defaultUsername),
				zap.String("password", defaultPassword),
			)
		}
	} else {
		fmt.Println("Default user already exists")
	}
}
