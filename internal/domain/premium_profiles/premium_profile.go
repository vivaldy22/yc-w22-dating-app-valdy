package premium_profiles

import (
	"time"
)

type PremiumProfile struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	UserID       string    `gorm:"type:uuid;not null;index"`
	Feature      string    `gorm:"type:varchar(50);not null"`
	PurchaseDate time.Time `gorm:"type:timestamptz;default:now()"`
	ExpiryDate   time.Time `gorm:"type:timestamptz;default:now()"`
}
