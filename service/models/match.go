package models

import "gorm.io/gorm"

type MatchResult struct {
	gorm.Model

	Players []User
	Scores  []int

	RoundCount int
}
