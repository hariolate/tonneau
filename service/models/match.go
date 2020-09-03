package models

import "gorm.io/gorm"

type MatchResult struct {
	gorm.Model

	Players []User `gorm:"foreignKey:ID"`
	Scores  []int

	RoundCount int
}
