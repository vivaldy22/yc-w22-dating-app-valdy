package postgres

import (
	"context"
	"errors"
	"log"
	"strings"

	"gorm.io/gorm"

	"yc-w22-dating-app-valdy/internal/domain/profiles"
	"yc-w22-dating-app-valdy/pkg/constant"
	"yc-w22-dating-app-valdy/pkg/database"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type (
	ProfileRepository interface {
		GetDB() *database.Database
		WithTx(tx *gorm.DB) ProfileRepository
		Create(ctx context.Context, profile *profiles.Profile) error
		FindByUserID(ctx context.Context, userID string) (profiles.Profile, error)
		FindSwipeableProfiles(ctx context.Context, swiperID, swiperGender string) ([]profiles.Profile, error)
	}

	profileRepository struct {
		db        *database.Database
		tableName string
	}
)

func NewProfileRepository(db *database.Database) ProfileRepository {
	if db == nil {
		panic("db is nil")
	}

	return &profileRepository{
		db:        db,
		tableName: profiles.TableName,
	}
}

func (pr *profileRepository) GetDB() *database.Database {
	return pr.db
}

func (pr *profileRepository) WithTx(tx *gorm.DB) ProfileRepository {
	return &profileRepository{
		db:        &database.Database{Master: tx},
		tableName: pr.tableName,
	}
}

func (pr *profileRepository) Create(ctx context.Context, profile *profiles.Profile) error {
	err := pr.db.Master.Create(profile).Error
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), constant.DuplicateRecord) {
		log.Printf("profile %s already exists", profile.ID)
		return ierror.ErrDuplicateData
	}

	log.Printf("Failed Create Profile: %s\n", err.Error())
	return ierror.ErrDatabase
}

func (pr *profileRepository) FindByUserID(ctx context.Context, userID string) (profiles.Profile, error) {
	entity := profiles.Profile{}
	err := pr.db.Master.First(&entity, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Profile with userID %s not found\n", userID)
		return entity, errors.New("data not found")
	}
	if err != nil {
		log.Printf("Failed Find Profile: %s\n", err.Error())
		return entity, errors.New("database error")
	}

	return entity, nil
}

func (pr *profileRepository) FindSwipeableProfiles(ctx context.Context, swiperID, swiperGender string) ([]profiles.Profile, error) {
	var entities []profiles.Profile

	err := pr.db.Slave.Raw(`
		SELECT p.*
		FROM profiles p
		LEFT JOIN swipes s
		  ON s.swiped_id = p.user_id AND s.swiper_id = ?
		WHERE s.id IS NULL
		  AND p.user_id != ?
		  AND p.gender != ?
		LIMIT 10;`, swiperID, swiperID, swiperGender).Scan(&entities).Error
	if len(entities) == 0 {
		log.Printf("No profiles found\n")
		return entities, ierror.ErrDataNotFound
	}
	if err != nil {
		log.Printf("Failed Find Profile: %s\n", err.Error())
		return entities, ierror.ErrDatabase
	}

	return entities, nil
}
