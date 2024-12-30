package swipes

import (
	"time"
)

const (
	TableName = "swipes"
)

type Swipe struct {
	ID       string    `gorm:"type:uuid;primaryKey"`
	SwiperID string    `gorm:"type:uuid;not null;uniqueIndex:unique_swiper_swiped"`
	SwipedID string    `gorm:"type:uuid;not null;uniqueIndex:unique_swiper_swiped"`
	Action   string    `gorm:"type:varchar(10);not null"`
	SwipedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (s Swipe) TableName() string {
	return TableName
}
