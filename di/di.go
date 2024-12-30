package di

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"yc-w22-dating-app-valdy/config"
	"yc-w22-dating-app-valdy/internal/repository/postgres"
	redisrepo "yc-w22-dating-app-valdy/internal/repository/redis"
	"yc-w22-dating-app-valdy/internal/usecase/auth"
	"yc-w22-dating-app-valdy/internal/usecase/swipe"
	"yc-w22-dating-app-valdy/pkg/database"
)

type DI struct {
	Configuration *config.Configuration
	Logger        *zap.Logger
	Echo          *echo.Echo
	Database      *database.Database

	AuthService  auth.Service
	SwipeService swipe.Service
}

func SetupDependencies() *DI {
	cfg := config.LoadConfig()
	e := echo.New()
	logger := zap.New(nil)
	db := database.NewDatabase(&cfg)
	redisClient := redis.NewClient(&cfg.RedisOption)

	// Setup Postgres Repositories
	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)

	// Setup Redis Repositories
	rateLimitRepo := redisrepo.NewRateLimitRepository(redisClient)

	// Setup Use cases
	authService := auth.NewService(&cfg, rateLimitRepo, userRepo, profileRepo)
	swipeService := swipe.NewService(rateLimitRepo, profileRepo)

	return &DI{
		Configuration: &cfg,
		Logger:        logger,
		Echo:          e,
		Database:      db,
		AuthService:   authService,
		SwipeService:  swipeService,
	}
}

func (d *DI) CleanUp() {
	log.Println("Cleaning up...")

	d.Database.CleanUp()
}
