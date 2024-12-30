package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joomcode/errorx"
	"gorm.io/gorm"

	"yc-w22-dating-app-valdy/config"
	"yc-w22-dating-app-valdy/internal/domain/profiles"
	"yc-w22-dating-app-valdy/internal/domain/users"
	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/internal/repository/postgres"
	"yc-w22-dating-app-valdy/internal/repository/redis"
	"yc-w22-dating-app-valdy/pkg/crypto"
	ierror "yc-w22-dating-app-valdy/pkg/error"
	ijwt "yc-w22-dating-app-valdy/pkg/jwt"
)

type (
	Service interface {
		SignUp(ctx context.Context, req model.SignUpRequest) (model.SignUpResponse, error)
		Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	}

	service struct {
		cfg           *config.Configuration
		rateLimitRepo redis.RateLimitRepository
		userRepo      postgres.UserRepository
		profileRepo   postgres.ProfileRepository
	}
)

func NewService(cfg *config.Configuration, rateLimitRepo redis.RateLimitRepository, userRepo postgres.UserRepository, profileRepo postgres.ProfileRepository) Service {
	if cfg == nil {
		panic("cfg is nil")
	}
	if rateLimitRepo == nil {
		panic("rateLimitRepo is nil")
	}
	if userRepo == nil {
		panic("userRepo is nil")
	}
	if profileRepo == nil {
		panic("profileRepo is nil")
	}

	return &service{
		cfg:           cfg,
		rateLimitRepo: rateLimitRepo,
		userRepo:      userRepo,
		profileRepo:   profileRepo,
	}
}

func (s *service) SignUp(ctx context.Context, req model.SignUpRequest) (model.SignUpResponse, error) {
	// Find existing user and profile
	fUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		_, err = s.profileRepo.FindByUserID(ctx, fUser.ID)
		if err == nil {
			log.Printf("profile with userID: %s already exists\n", fUser.ID)
			return model.SignUpResponse{}, errors.New("user already exists")
		}

		log.Printf("user with email: %s already exists\n", req.Email)
		return model.SignUpResponse{}, errors.New("user already exists")
	}

	// Database error
	if !errorx.IsNotFound(err) {
		return model.SignUpResponse{}, err
	}

	// Generate User ID using google uuid v7
	userID, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate user id: %s\n", err.Error())
		return model.SignUpResponse{}, ierror.ErrGeneral
	}

	// Generate Profile ID using google uuid v7
	profileID, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate profile id: %s\n", err.Error())
		return model.SignUpResponse{}, ierror.ErrGeneral
	}

	// Decrypt password from frontend
	plainPassword, err := crypto.Decrypt(req.Password, s.cfg.HashSecret)
	if err != nil {
		log.Printf("failed to decrypt password: %s\n", err.Error())
		return model.SignUpResponse{}, ierror.ErrGeneral
	}

	// Hash and Salt password before storing to db
	hashedPassword, err := crypto.HashAndSalt(plainPassword)
	if err != nil {
		log.Printf("failed to hash password: %s\n", err.Error())
		return model.SignUpResponse{}, ierror.ErrGeneral
	}

	user := users.User{
		ID:           userID.String(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	profile := profiles.Profile{
		ID:       profileID.String(),
		UserID:   userID.String(),
		Name:     req.Name,
		Bio:      req.Bio,
		Age:      req.Age,
		Gender:   req.Gender,
		Location: req.Location,
		PhotoURL: req.PhotoURL,
	}

	// Use DB Tx to Create User and Profile entity
	err = s.userRepo.GetDB().UseTx(ctx, func(tx *gorm.DB) error {
		if errTx := s.userRepo.WithTx(tx).Create(ctx, &user); errTx != nil {
			return errTx
		}

		if errTx := s.profileRepo.WithTx(tx).Create(ctx, &profile); errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return model.SignUpResponse{}, err
	}

	return model.SignUpResponse{
		Name: user.Name,
	}, nil
}

func (s *service) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return model.LoginResponse{}, err
	}

	// Decrypt password from frontend
	plainPassword, err := crypto.Decrypt(req.Password, s.cfg.HashSecret)
	if err != nil {
		log.Printf("failed to decrypt password: %s\n", err.Error())
		return model.LoginResponse{}, ierror.ErrGeneral
	}

	if !crypto.CheckPasswordHash(user.PasswordHash, plainPassword) {
		return model.LoginResponse{}, ierror.ErrInvalidPassword
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"id":          user.ID,
		"name":        user.Name,
		"gender":      user.Gender,
		"is_verified": user.IsVerified,
		"exp":         time.Now().Add(time.Hour * 24).UnixMilli(),
	}
	accessToken, err := ijwt.GenerateJWT(s.cfg.JWTSecret, claims)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		AccessToken: accessToken,
	}, nil

}
