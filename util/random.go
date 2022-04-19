package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// random generator for owner name
func RandomOwner() string {
	return RandomString(7)
}

// random money generator for balance
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// random currency generator
func RandomCurrency() string {
	currencies := []string{"USD", "JPY", "CAD", "EUR", "NZD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
