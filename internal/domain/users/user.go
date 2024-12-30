package users

import (
	"time"
)

const (
	TableName = "users"
)

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	Name         string    `gorm:"type:varchar(100);not null"`
	IsVerified   bool      `gorm:"type:boolean;default:false"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:now()"`
}

func (u User) TableName() string {
	return TableName
}
