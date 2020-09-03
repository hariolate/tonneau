package models

import "github.com/go-playground/validator"

type NewPassword struct {
	Password string `json:"password" validate:"required,min=3"`
}

func (n *NewPassword) Validate() error {
	return validator.New().Struct(n)
}
