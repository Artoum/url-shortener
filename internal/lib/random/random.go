package random

import (
	"math/rand/v2"
	"strings"
)

const allowLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func NewRandomString(aliasLength int) string {
	var builder strings.Builder

	builder.Grow(aliasLength)
	for i := 0; i < aliasLength; i++ {
		randomIndex := rand.IntN(len(allowLetters))
		builder.WriteByte(allowLetters[randomIndex])
	}

	return builder.String()
}
