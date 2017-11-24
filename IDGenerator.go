package main

import (
    "time"
    "math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
    rand.Seed(time.Now().UnixNano())
}
//
// Create new random charge ID Sequence
//
func CreateChargeId() string {
    return "1TESTAGQ2G0H1tnT" + GetRandomString(8)
}

//
// Create new random Cart ID Sequence
//
func CreateCardId() string {
    return "1TESTAGQ2G0H1tnT" + GetRandomString(8)
}

//
// Create new random finger print Sequence
//
func CreateFingerPrint() string {
    return GetRandomString(16)
}

//
// Create new random Request id Sequence
//
func CreateRequestId() string{
    //UxSQbUq877Jnyb
    return "req_" + GetRandomString(14)
}

//
// Generate new random string
//
func GetRandomString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}
