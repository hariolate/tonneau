package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model

	Alias    string        `json:"alias"`
	Matches  []MatchResult `json:"matches"`
	Picture  []byte        `json:"picture"`
	Trophies uint          `json:"trophies"`
}
