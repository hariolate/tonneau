package models

import (
	"github.com/go-playground/validator"
)

type Signup struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

func (s *Signup) Validate() error {
	return validator.New().Struct(s)
}
