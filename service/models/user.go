package models

import (
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique" validate:"required,email"`
	Password string `validate:"required,min=3"`

	Profile Profile
}

func (u *User) Validate() error {
	return validator.New().Struct(u)
}
