package swipe

import (
	"context"

	"yc-w22-dating-app-valdy/internal/model"
	"yc-w22-dating-app-valdy/internal/repository/postgres"
	"yc-w22-dating-app-valdy/internal/repository/redis"
)

type (
	Service interface {
		GetSwipeableProfiles(ctx context.Context, req model.GetSwipeableProfilesRequest) (model.GetSwipeableProfilesResponse, error)
	}

	service struct {
		rateLimitRepo redis.RateLimitRepository
		profileRepo   postgres.ProfileRepository
	}
)

func NewService(rateLimitRepo redis.RateLimitRepository, profileRepo postgres.ProfileRepository) Service {
	if rateLimitRepo == nil {
		panic("rateLimitRepo is nil")
	}
	if profileRepo == nil {
		panic("profileRepo is nil")
	}

	return &service{
		rateLimitRepo: rateLimitRepo,
		profileRepo:   profileRepo,
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
