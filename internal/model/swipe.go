package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	SwipeRequest struct {
		SwiperID string `json:"-"`
		SwipedID string `json:"swiped_id"`
	}

	SwipeResponse struct {
		IsMutual     bool  `json:"is_mutual"`
		DailyCounter int64 `json:"daily_counter"`
	}
)

func (m SwipeRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SwiperID, validation.Required),
		validation.Field(&m.SwipedID, validation.Required),
	)
}
