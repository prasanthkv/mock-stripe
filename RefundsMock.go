package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"fmt"
)
func RefundsHandler(w http.ResponseWriter, r *http.Request) {
	//all request are json
	w.Header().Set("content-type", "application/json")
	w.Header().Set("stripe-version", "mock-1.0")
	//capture id
	vars := mux.Vars(r)
	captureId := vars["id"]
	//copy the idempotency key to response
	idempotencyKey := r.Header.Get("idempotency-key");
	w.Header().Set("idempotency-key", idempotencyKey)
	log.Println(w, "{\"chargeId\":\""+captureId+"\"}:")
	fmt.Fprintln(w, "{\"chargeId\":\""+captureId+"\"}:")
	//should be the last
	w.WriteHeader(http.StatusOK)
}