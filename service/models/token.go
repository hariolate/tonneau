package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Token struct {
	UID   uint
	Token uuid.UUID
}

func NewToken(UID uint, token uuid.UUID) *Token {
	return &Token{UID: UID, Token: token}
}

func NewTokenFor(user *User) *Token {
	return &Token{
		user.ID,
		uuid.New(),
	}
}

func (t *Token) RedisKey() string {
	return fmt.Sprintf("user:%d:tokens", t.UID)
}

func (t *Token) RedisToken() string {
	return t.Token.String()
}
