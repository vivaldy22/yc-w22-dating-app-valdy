package model

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Bio      string `json:"bio"`
		Age      int32  `json:"age"`
		Gender   string `json:"gender"`
		Location string `json:"location"`
		PhotoURL string `json:"photo_url"`
	}

	SignUpResponse struct {
		Name string `json:"name"`
	}
)

func (m *SignUpRequest) Validate() error {
	m.Gender = strings.ToLower(m.Gender)

	return validation.ValidateStruct(m,
		validation.Field(&m.Email, validation.Required, is.Email, validation.Length(0, 255)),
		validation.Field(&m.Name, validation.Required, validation.Length(0, 100)),
		validation.Field(&m.Bio),
		validation.Field(&m.Age, validation.Required, validation.Min(18)),
		validation.Field(&m.Gender, validation.Required, validation.In("m", "f")),
		validation.Field(&m.Location, validation.Required, validation.Length(0, 100)),
		validation.Field(&m.PhotoURL),
	)
}
