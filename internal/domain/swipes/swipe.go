package swipes

import (
	"time"
)

type Swipe struct {
	ID       string    `gorm:"type:uuid;primaryKey"`
	SwiperID string    `gorm:"type:uuid;not null"`
	SwipedID string    `gorm:"type:uuid;not null"`
	Action   string    `gorm:"type:varchar(10);not null"`
	SwipedAt time.Time `gorm:"type:timestamptz;default:now()"`
}
