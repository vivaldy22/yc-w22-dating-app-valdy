package di

import (
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"yc-w22-dating-app-valdy/config"
	"yc-w22-dating-app-valdy/internal/repository/postgres"
	redisrepo "yc-w22-dating-app-valdy/internal/repository/redis"
	"yc-w22-dating-app-valdy/internal/usecase/auth"
	"yc-w22-dating-app-valdy/internal/usecase/onboard"
	"yc-w22-dating-app-valdy/pkg/database"
	"yc-w22-dating-app-valdy/pkg/redis"
)

type DI struct {
	Configuration *config.Configuration
	Logger        *zap.Logger
	Echo          *echo.Echo
	Database      *database.Database

	AuthService    auth.Service
	OnboardService onboard.Service
}

func SetupDependencies() *DI {
	cfg := config.LoadConfig()
	e := echo.New()
	logger := zap.New(nil)
	db := database.NewDatabase(&cfg)
	redisClient := redis.NewRedis(&cfg.Redis)

	// Setup Postgres Repositories
	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	swipeRepo := postgres.NewSwipeRepository(db)
	premiumProfileRepo := postgres.NewPremiumProfileRepository(db)

	// Setup Redis Repositories
	rateLimitRepo := redisrepo.NewRateLimitRepository(redisClient)

	// Setup Use cases
	authService := auth.NewService(&cfg, rateLimitRepo, userRepo, profileRepo)
	onboardService := onboard.NewService(&cfg, rateLimitRepo, userRepo, profileRepo, swipeRepo, premiumProfileRepo)

	return &DI{
		Configuration:  &cfg,
		Logger:         logger,
		Echo:           e,
		Database:       db,
		AuthService:    authService,
		OnboardService: onboardService,
	}
}

func (d *DI) CleanUp() {
	log.Println("Cleaning up...")

	d.Database.CleanUp()
}
