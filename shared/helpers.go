package shared

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

func RandomInRange(from, to int) int {
	return from + rand.Int()%(to-from)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	return bcrypt.CompareHashAndPassword(byteHash, plainPwd) == nil
}
