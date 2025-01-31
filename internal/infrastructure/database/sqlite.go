package database

import (
	"fmt"

	"github.com/mrfansi/go-api-boilerplate/internal/domain/entity"
	"github.com/mrfansi/go-api-boilerplate/internal/infrastructure/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable foreign key constraints
	db.Exec("PRAGMA foreign_keys = ON")

	// Auto migrate the schema
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate schema: %w", err)
	}

	return db, nil
}

// autoMigrate automatically migrates the schema for all entities
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		// Add other entities here as they are created
	)
}

// Transaction executes the given function within a database transaction
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
