package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	GetSwipeableProfilesRequest struct {
		UserID string `json:"user_id"`
		Gender string `json:"gender"`
	}

	SwipeableProfile struct {
		UserID   string `json:"user_id"`
		Name     string `json:"name"`
		Bio      string `json:"bio"`
		Age      int32  `json:"age"`
		Gender   string `json:"gender"`
		Location string `json:"location"`
		PhotoURL string `json:"photo_url"`
	}

	GetSwipeableProfilesResponse struct {
		Profiles []SwipeableProfile `json:"profiles"`
	}
)

func (m GetSwipeableProfilesRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
		validation.Field(&m.Gender, validation.Required, validation.In("m", "f")),
	)
}
