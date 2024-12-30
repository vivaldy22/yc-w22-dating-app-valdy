package database

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"yc-w22-dating-app-valdy/config"
	"yc-w22-dating-app-valdy/internal/domain/premium_profiles"
	"yc-w22-dating-app-valdy/internal/domain/profiles"
	"yc-w22-dating-app-valdy/internal/domain/swipes"
	"yc-w22-dating-app-valdy/internal/domain/users"
)

type Database struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func NewDatabase(cfg *config.Configuration) *Database {
	log.Println("Starting Master Database Connection")
	master := newDatabase(&cfg.DatabaseMaster, cfg.FeatureFlag.EnableDatabaseAutoMigrate)
	log.Println("Starting Slave Database Connection")
	slave := newDatabase(&cfg.DatabaseSlave, cfg.FeatureFlag.EnableDatabaseAutoMigrate)

	return &Database{
		Master: master,
		Slave:  slave,
	}
}

func newDatabase(cfg *config.Database, isAutoMigrate bool) *gorm.DB {
	dbConnectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable search_path=%s",
		cfg.Username,
		cfg.Password,
		cfg.Name,
		cfg.Host,
		cfg.Port,
		cfg.Schema,
	)
	gormOpts := gorm.Config{}

	db, err := gorm.Open(postgres.Open(dbConnectionString), &gormOpts)
	if err != nil {
		panic(err)
	}

	if isAutoMigrate {
		log.Printf("Starting Auto-Migrate Database %s in Schema %s\n", cfg.Name, cfg.Schema)
		err = db.AutoMigrate(
			&users.User{},
			&profiles.Profile{},
			&swipes.Swipe{},
			&premium_profiles.PremiumProfile{},
		)
		if err != nil {
			panic(err)
		}
		log.Printf("Finished Auto-Migrate Database %s in Schema %s\n", cfg.Name, cfg.Schema)
	}

	return db
}

func (d Database) UseTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	var (
		tx  = d.Master.Begin()
		err error
	)

	defer func() {
		if err != nil {
			errRb := tx.Rollback()
			log.Printf("TxRollback err: %v\n", errRb)
		}
	}()

	if err = tx.Error; err == nil {
		err = fn(tx)
	}

	if err == nil {
		err = tx.Commit().Error
		log.Printf("TxCommit err: %v\n", err)
	}

	return err
}

func (d Database) CleanUp() {
	log.Printf("CleanUp Database...")

	master, err := d.Master.DB()
	if err != nil {
		log.Printf("CleanUp Master Database failed: %s\n", err.Error())
	}
	slave, err := d.Slave.DB()
	if err != nil {
		log.Printf("CleanUp Slave Database failed: %s\n", err.Error())
	}
	err = master.Close()
	if err != nil {
		log.Printf("CleanUp Master Database failed: %s\n", err.Error())
	}
	err = slave.Close()
	if err != nil {
		log.Printf("CleanUp Slave Database failed: %s\n", err.Error())
	}
}
