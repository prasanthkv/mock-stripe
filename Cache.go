package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

//--
var authCache = cache.New(5*time.Hour, 5*time.Hour)
var chargeCache = cache.New(5*time.Hour, 5*time.Hour)
var idempotencyCache = cache.New(5*time.Hour, 5*time.Hour)

type CacheObject struct {
	Status      int
	RequestHash string
	Idempotency string
	Charge      ChargeObject
	Error       ErrorObject
}
