package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (m LoginRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, is.Email),
		validation.Field(&m.Password, validation.Required),
	)
}
