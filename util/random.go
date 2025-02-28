package util

import (
	"math/rand/v2"
	"strings"
)

const alphanumericSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphanumericSymbols)

	for i := 0; i < n; i++ {
		c := alphanumericSymbols[rand.IntN(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
