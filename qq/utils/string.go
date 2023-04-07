package utils

import (
	"math/rand"
	"strings"
	"time"
)

var utilRand *rand.Rand

func init() {
	utilRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(len int) string {
	return RandomStringRange(len, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func RandomStringRange(length int, str string) string {
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(str[utilRand.Intn(len(str))])
	}
	return sb.String()
}
