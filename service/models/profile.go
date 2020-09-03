package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Alias string `json:"alias"`
	//Matches  []MatchResult `json:"matches" gorm:"foreignKey:ID"`
	Picture  []byte `json:"picture"`
	Trophies int    `json:"trophies"`
	UserID   uint
}
