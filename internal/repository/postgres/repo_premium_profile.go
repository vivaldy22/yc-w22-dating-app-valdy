package postgres

import (
	"context"
	"log"

	"gorm.io/gorm"

	"yc-w22-dating-app-valdy/internal/domain/premium_profiles"
	"yc-w22-dating-app-valdy/internal/domain/profiles"
	"yc-w22-dating-app-valdy/pkg/database"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type (
	PremiumProfileRepository interface {
		GetDB() *database.Database
		WithTx(tx *gorm.DB) PremiumProfileRepository
		Create(ctx context.Context, premium *premium_profiles.PremiumProfile) error
	}

	premiumProfileRepository struct {
		db        *database.Database
		tableName string
	}
)

func NewPremiumProfileRepository(db *database.Database) PremiumProfileRepository {
	if db == nil {
		panic("db is nil")
	}

	return &premiumProfileRepository{
		db:        db,
		tableName: profiles.TableName,
	}
}

func (pr *premiumProfileRepository) GetDB() *database.Database {
	return pr.db
}

func (pr *premiumProfileRepository) WithTx(tx *gorm.DB) PremiumProfileRepository {
	return &premiumProfileRepository{
		db:        &database.Database{Master: tx},
		tableName: pr.tableName,
	}
}

func (pr *premiumProfileRepository) Create(ctx context.Context, premium *premium_profiles.PremiumProfile) error {
	err := pr.db.Master.Create(premium).Error
	if err == nil {
		return nil
	}

	log.Printf("Failed Create Profile: %s\n", err.Error())
	return ierror.ErrDatabase
}
