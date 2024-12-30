package postgres

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"yc-w22-dating-app-valdy/internal/domain/users"
	"yc-w22-dating-app-valdy/pkg/database"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type UserRepository interface {
	GetDB() *database.Database
	WithTx(tx *gorm.DB) UserRepository
	Create(ctx context.Context, user *users.User) error
	FindByEmail(ctx context.Context, email string) (users.User, error)
}

type userRepository struct {
	db        *database.Database
	tableName string
}

func NewUserRepository(db *database.Database) UserRepository {
	if db == nil {
		panic("db is nil")
	}

	return &userRepository{
		db:        db,
		tableName: users.TableName,
	}
}

func (ur *userRepository) GetDB() *database.Database {
	return ur.db
}

func (ur *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{
		db:        &database.Database{Master: tx},
		tableName: ur.tableName,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *users.User) error {
	err := ur.db.Master.Create(user).Error
	if err != nil {
		log.Printf("Failed Create User: %s\n", err.Error())
		return ierror.ErrDatabase
	}

	return nil
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (users.User, error) {
	entity := users.User{}
	err := ur.db.Slave.First(&entity, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("User with email %s not found\n", email)
		return entity, ierror.ErrDataNotFound
	}
	if err != nil {
		log.Printf("Failed Find User: %s\n", err.Error())
		return entity, ierror.ErrDatabase
	}

	return entity, nil
}
