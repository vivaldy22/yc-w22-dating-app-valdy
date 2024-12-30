package onboard

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joomcode/errorx"
	"gorm.io/gorm"

	"yc-w22-dating-app-valdy/config"
	"yc-w22-dating-app-valdy/internal/domain/premium_profiles"
	"yc-w22-dating-app-valdy/internal/domain/swipes"
	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/internal/repository/postgres"
	"yc-w22-dating-app-valdy/internal/repository/redis"
	ierror "yc-w22-dating-app-valdy/pkg/error"
)

type (
	Service interface {
		GetSwipeableProfiles(ctx context.Context, req model.GetSwipeableProfilesRequest) (model.GetSwipeableProfilesResponse, error)
		Swipe(ctx context.Context, req model.SwipeRequest, action string) (model.SwipeResponse, error)
		BuyPremiumFeature(ctx context.Context, req model.BuyPremiumFeatureRequest) (model.BuyPremiumFeatureResponse, error)
	}

	service struct {
		cfg                *config.Configuration
		rateLimitRepo      redis.RateLimitRepository
		userRepo           postgres.UserRepository
		profileRepo        postgres.ProfileRepository
		swipeRepo          postgres.SwipeRepository
		premiumProfileRepo postgres.PremiumProfileRepository
	}
)

func NewService(cfg *config.Configuration, rateLimitRepo redis.RateLimitRepository, userRepo postgres.UserRepository, profileRepo postgres.ProfileRepository, swipeRepo postgres.SwipeRepository, premiumProfileRepo postgres.PremiumProfileRepository) Service {
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
	if swipeRepo == nil {
		panic("swipeRepo is nil")
	}
	if premiumProfileRepo == nil {
		panic("premiumProfileRepo is nil")
	}

	return &service{
		cfg:                cfg,
		rateLimitRepo:      rateLimitRepo,
		userRepo:           userRepo,
		profileRepo:        profileRepo,
		swipeRepo:          swipeRepo,
		premiumProfileRepo: premiumProfileRepo,
	}
}

func (s *service) GetSwipeableProfiles(ctx context.Context, req model.GetSwipeableProfilesRequest) (model.GetSwipeableProfilesResponse, error) {
	profiles, err := s.profileRepo.FindSwipeableProfiles(ctx, req.UserID, req.Gender)
	if err != nil {
		return model.GetSwipeableProfilesResponse{}, err
	}

	var profilesResponse []model.SwipeableProfile
	for _, profile := range profiles {
		profilesResponse = append(profilesResponse, model.SwipeableProfile{
			UserID:   profile.UserID,
			Name:     profile.Name,
			Bio:      profile.Bio,
			Age:      profile.Age,
			Gender:   profile.Gender,
			Location: profile.Location,
			PhotoURL: profile.PhotoURL,
		})
	}

	return model.GetSwipeableProfilesResponse{
		Profiles: profilesResponse,
	}, nil
}

func (s *service) Swipe(ctx context.Context, req model.SwipeRequest, action string) (model.SwipeResponse, error) {
	res := model.SwipeResponse{}

	// Check Rate Limit
	counter, err := s.rateLimitRepo.Get(ctx, "onboard#daily", req.SwiperID)
	if err != nil && !errorx.IsNotFound(err) {
		return model.SwipeResponse{}, ierror.ErrGeneral
	}

	res.DailyCounter = counter

	if req.SwipedID == req.SwiperID {
		return res, ierror.ErrNoSelfSwipe
	}

	if s.cfg.SwipeDailyLimit > 0 && counter >= s.cfg.SwipeDailyLimit {
		return res, ierror.ErrDailySwipeLimitReached
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate onboard id: %s\n", err.Error())
		return res, ierror.ErrGeneral
	}

	swipe := swipes.Swipe{
		ID:       id.String(),
		SwiperID: req.SwiperID,
		SwipedID: req.SwipedID,
		Action:   action,
	}

	err = s.swipeRepo.Create(ctx, &swipe)
	if errorx.IsDuplicate(err) {
		return res, ierror.ErrProfileAlreadySwiped
	}
	if err != nil {
		return res, err
	}

	isMutual, err := s.swipeRepo.CheckMutualLike(ctx, req.SwiperID, req.SwipedID)
	if err != nil {
		return res, err
	}

	counter, err = s.rateLimitRepo.Incr(ctx, "onboard#daily", req.SwiperID, 24*time.Hour)
	if err != nil {
		return res, err
	}

	return model.SwipeResponse{
		IsMutual:     isMutual,
		DailyCounter: counter,
	}, nil
}

func (s *service) BuyPremiumFeature(ctx context.Context, req model.BuyPremiumFeatureRequest) (model.BuyPremiumFeatureResponse, error) {
	res := model.BuyPremiumFeatureResponse{
		Feature: req.Feature,
	}
	purchaseDate := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate premium id: %s\n", err.Error())
		return res, ierror.ErrGeneral
	}

	err = s.premiumProfileRepo.GetDB().UseTx(ctx, func(tx *gorm.DB) error {
		errTx := s.premiumProfileRepo.WithTx(tx).Create(ctx, &premium_profiles.PremiumProfile{
			ID:           id.String(),
			UserID:       req.UserID,
			Feature:      req.Feature,
			PurchaseDate: purchaseDate,
			ExpiryDate:   purchaseDate.Add(30 * 24 * time.Hour),
		})
		if errTx != nil {
			return errTx
		}

		switch req.Feature {
		case "verified_user":
			errTx = s.userRepo.WithTx(tx).UpdateVerified(ctx, req.UserID)
			if errTx != nil {
				return errTx
			}
		default:
			log.Println("feature not available to purchase")
			return ierror.ErrInvalidRequest
		}

		return nil
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
