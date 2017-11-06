package main

import (
	"time"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/patrickmn/go-cache"
)

//--
var authCache = cache.New(5*time.Hour, 5*time.Hour)
var chargeCache = cache.New(5*time.Hour, 5*time.Hour)
var idempotencyCache = cache.New(5*time.Hour, 5*time.Hour)

type CacheObject struct {
	Status      int
	RequestId   string
	RequestHash string
	Idempotency string
	Charge      ChargeObject
	Error       ErrorResponse
	Type		string
}

func MD5Hash(v interface{}) (string) {
	out, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	aStringToHash := []byte(out)
	//Get the hashes in bytes
	md5Bytes := md5.Sum(aStringToHash)
	//hash
	return hex.EncodeToString(md5Bytes[:])
}
