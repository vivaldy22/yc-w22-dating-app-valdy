package postgres

import (
	"context"
	"log"
	"strings"

	"yc-w22-dating-app-valdy/internal/domain/swipes"
	"yc-w22-dating-app-valdy/pkg/constant"
	"yc-w22-dating-app-valdy/pkg/database"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type (
	SwipeRepository interface {
		Create(ctx context.Context, swipe *swipes.Swipe) error
		CheckMutualLike(ctx context.Context, swiperID, swipedID string) (bool, error)
	}

	swipeRepository struct {
		db        *database.Database
		tableName string
	}
)

func NewSwipeRepository(db *database.Database) SwipeRepository {
	if db == nil {
		panic("db is nil")
	}

	return &swipeRepository{
		db:        db,
		tableName: swipes.TableName,
	}
}

func (sr *swipeRepository) Create(ctx context.Context, swipe *swipes.Swipe) error {
	err := sr.db.Master.Create(swipe).Error
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), constant.DuplicateRecord) {
		log.Printf("onboard %s already exists", swipe.ID)
		return ierror.ErrDuplicateData
	}

	log.Printf("Failed Create Swipe: %s\n", err.Error())
	return ierror.ErrDatabase
}

func (sr *swipeRepository) CheckMutualLike(ctx context.Context, swiperID, swipedID string) (bool, error) {
	exists := false
	err := sr.db.Slave.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM swipes AS s1
			JOIN swipes AS s2
			  ON s1.swiper_id = s2.swiped_id
			 AND s1.swiped_id = s2.swiper_id
			WHERE s1.swiper_id = ?
			  AND s1.swiped_id = ?
			  AND s1.action = 'like'
			  AND s2.action = 'like'
		);`, swiperID, swipedID).Scan(&exists).Error
	if err != nil {
		return false, ierror.ErrDatabase
	}

	return exists, nil
}
