package util

import (
	"math/rand"
	"strconv"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreateTransaction() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000000000000
	max := 999999999999999

	randomNumber := rand.Intn(max-min+1) + min
	randomString := strconv.Itoa(randomNumber)
	return randomString
}
