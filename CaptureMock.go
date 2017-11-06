package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"encoding/json"
)

func CaptureHandler(w http.ResponseWriter, r *http.Request) {
	//all request are json
	w.Header().Set("content-type", "application/json")
	w.Header().Set("stripe-version", "mock-1.0")
	//capture id
	vars := mux.Vars(r)
	captureId := vars["id"]
	//copy the idempotency key to response
	idempotencyKey := r.Header.Get("idempotency-key");
	w.Header().Set("idempotency-key", idempotencyKey)
	//evaluate idempotency key
	cObject, found := authCache.Get(idempotencyKey)
	if (found) {
		print("worng idempotency-key", cObject.(string))
		errorObjects := ErrorResponse{
			Error: ErrorObject{
				Type:    "idempotency_error",
				Message: "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges/" + captureId + "/capture' vs '/v1/charges'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request.",
			},
		}
		fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObjects))
		//should be the last
		w.WriteHeader(http.StatusBadRequest)
	} else {
		//
		log.Println(w, "{\"chargeId\":\""+captureId+"\"}:")
		//return
		//should be the last
		w.WriteHeader(http.StatusOK)
	}
}
