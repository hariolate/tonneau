package models

import "github.com/go-playground/validator"

type NewEmail struct {
	Email string `json:"email" validate:"required,email"`
}

func (e *NewEmail) Validate() error {
	return validator.New().Struct(e)
}
