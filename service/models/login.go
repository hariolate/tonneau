package models

import "github.com/go-playground/validator"

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

func (l *Login) Validate() error {
	return validator.New().Struct(l)
}
