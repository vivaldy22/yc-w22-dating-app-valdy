package profiles

import (
	"time"
)

const (
	TableName = "profiles"
)

type Profile struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:varchar(100)"`
	Bio       string    `gorm:"type:text"`
	Age       int32     `gorm:"type:integer"`
	Gender    string    `gorm:"type:varchar(1)"`
	Location  string    `gorm:"type:varchar(100)"`
	PhotoURL  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (p Profile) TableName() string {
	return TableName
}
