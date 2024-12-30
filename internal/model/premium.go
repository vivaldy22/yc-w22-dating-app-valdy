package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	BuyPremiumFeatureRequest struct {
		UserID  string `json:"-"`
		Feature string `json:"feature"`
	}

	BuyPremiumFeatureResponse struct {
		Feature string `json:"feature"`
	}
)

func (m BuyPremiumFeatureRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
		validation.Field(&m.Feature, validation.Required),
	)
}
