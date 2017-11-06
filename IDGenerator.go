package main

import (
	"time"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateChargeId() string {
	return "1BKXYAGQ2G0H1tnT" + GetRandomString(8)
}

func CreateCardId() string {
	return "1BKXYAGQ2G0H1tnT" + GetRandomString(8)
}

func CreateFingerPrint() string {
	return GetRandomString(16)
}

func GetRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateRequestId() string{
	//UxSQbUq877Jnyb
	return "req_" + GetRandomString(14)
}

type Todo struct {
	ID string `json:"id"`
}
