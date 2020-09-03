package models

type Profile struct {
	Alias    string        `json:"alias"`
	Matches  []MatchResult `json:"matches"`
	Picture  []byte        `json:"picture"`
	Trophies uint          `json:"trophies"`
}
