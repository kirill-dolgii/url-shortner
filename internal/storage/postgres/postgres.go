package postgres

import (
	"errors"
	"fmt"

	"github.com/kirill-dolgii/url-shortner/internal/config/dbconfig"
	"github.com/kirill-dolgii/url-shortner/internal/domain/models"
	"github.com/kirill-dolgii/url-shortner/internal/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func InitDB(cfg *dbconfig.DBConfig) (*Storage, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(url, alias string) (int64, error) {
	const op = "storage.postgres.Storage.SaveURL"
	model := models.Url{
		FullUrl: url,
		Alias:   alias,
	}

	if err := s.db.Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, storage.ErrUrlExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return model.ID, nil
}

func (s *Storage) GetUrlByAlias(alias string) (models.Url, error) {
	const op = "storage.postgres.Storage.GetUrlByAlias"
	var url models.Url
	err := s.db.First(&url, "alias = ?", alias).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Url{}, storage.ErrUrlNotFound
		}
		return models.Url{}, fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}
