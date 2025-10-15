package db

import (
	"fmt"
	"task-management/internal/config"
	"task-management/internal/infra/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Database struct {
	DB *sqlx.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sqlx.Connect(cfg.Driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenCons)
	db.SetMaxIdleConns(cfg.MaxIdleCons)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Minute)

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.String("name", cfg.Name),
	)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	if d.DB != nil {
		if err := d.DB.Close(); err != nil {
			logger.Error("failed to close database", zap.Error(err))
			return err
		}

		logger.Info("database connection closed")
	}

	return nil
}
