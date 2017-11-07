package main

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

func CaptureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CaptureHandler : INIT")
	//build the form
	r.ParseForm()
	fmt.Println("-----")
	fmt.Println(r.Form)
	fmt.Println("-----")
	httpStatus := http.StatusBadRequest
	//capture id
	vars := mux.Vars(r)
	captureId := vars["id"]
	//
	// Set all Headers
	//
	header := w.Header()
	requestId := CreateRequestId()
	header.Set("content-type", "application/json")
	header.Set("stripe-version", "mock-1.0")
	header.Set("request-id", requestId)
	header.Set("original-capture-id", captureId)
	//copy the idempotency key to response
	idempotencyKey := r.Header.Get("idempotency-key")
	header.Set("idempotency-key", idempotencyKey)
	//make hash of form this will help maintain idempotency
	formHash := MD5Hash(r.Form)
	header.Set("request-md5", formHash)
	//
	// check for auth idempotency_key
	//
	idempotencyObj, found := idempotencyCache.Get(idempotencyKey)
	//found
	if found {
		fmt.Println("CaptureHandler : idempotency found for key: " + idempotencyKey)
		//response object
		errorObjects := ErrorResponse{
			Error: ErrorObject{
				Type: "idempotency_error",
			},
		}
		//
		idempotency := idempotencyObj.(Idempotency)
		//
		exit := true
		if idempotency.Type == "auth" {
			fmt.Println("CaptureHandler : idempotency auth key found")
			//end user is trying to access capture with same idempotency as of auth.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges/" + captureId + "/refunds' vs '/v1/charges'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else if idempotency.Type == "capture" {
			fmt.Println("CaptureHandler : idempotency capture key found")
			//end user is trying to access capture with same idempotency as of capture.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges/" + captureId + "/refunds' vs '/v1/charges/" + captureId + "/capture'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else if idempotency.Type == "void" && idempotency.RequestHash != formHash {
			fmt.Println("CaptureHandler : idempotency void key found with md5:" + formHash)
			//end user is trying to access capture with same idempotency as of void with different form parameters.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used with the same parameters they were first used with. Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else {
			//valid request lets process
			fmt.Println("CaptureHandler : idempotency auth key cache")
			exit = false
		}
		// idempotency error so exit
		if exit {
			fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObjects))
			//final http status code
			w.WriteHeader(httpStatus)
			return
		}
	} else {
		//new request
		fmt.Println("CaptureHandler : new request with idempotency key: " + idempotencyKey)
		idempotency := Idempotency{
			Type:        "void",
			RequestId:   requestId,
			ChargeId:    captureId,
			RequestHash: formHash,
		}
		//set cache for next use
		idempotencyCache.Set(idempotencyKey, idempotency, cache.DefaultExpiration)
	}
	//
	//check for cached void object
	//
	cachedObj, found := captureCache.Get(captureId)
	if found {
		fmt.Println("CaptureHandler : cache fault:" + idempotencyKey)
		//get capture object from cache
		cacheObject := cachedObj.(CacheObject)
		//copy original request id
		header.Set("original-request", cacheObject.RequestId)
		//write to stream
		if cacheObject.Status == 200 {
			json.NewEncoder(w).Encode(cacheObject.Charge)
		} else {
			//this will never happen
			json.NewEncoder(w).Encode(cacheObject.Error)
		}
		//should be the last
		w.WriteHeader(cacheObject.Status)
		return
	}
	//
	// First time request
	//
	fmt.Println("CaptureHandler :First time request")
	//original request id and request id will be same this case
	header.Set("original-request", requestId)
	//evaluate auth
	cachedObj, found = chargeCache.Get(captureId)
	//
	// all set
	//
	if found {
		fmt.Println("CaptureHandler : Auth found for " + captureId)
		//process
		chargeObject := (cachedObj.(CacheObject)).Charge
		reqAmount, err := strconv.Atoi(FindFist(r.Form["amount"]))
		//
		fmt.Println("CaptureHandler : Start", err)
		//
		if reqAmount <= 0 {
			//charge amount should not be less than requested amount
			errorObjects := ErrorResponse{
				Error: ErrorObject{
					Type:    "invalid_request_error",
					Message: "Amount must be at least 50 cents.",
				},
			}
			//write to stream
			json.NewEncoder(w).Encode(errorObjects)
		} else if chargeObject.Amount < reqAmount {
			//charge amount should not be less than requested amount
			errorObjects := ErrorResponse{
				Error: ErrorObject{
					Type:    "invalid_request_error",
					Message: "You cannot capture a charge for an amount greater than it already has.",
				},
			}
			//write to stream
			json.NewEncoder(w).Encode(errorObjects)
		} else if chargeObject.Captured {
			//all ready captured
			json.NewEncoder(w).Encode(chargeObject)
			httpStatus = http.StatusOK
		} else {
			refundAmount := chargeObject.Amount - reqAmount
			refundId := "txn_" + CreateChargeId()
			chargeObject.Captured = true
			chargeObject.AmountRefunded = refundAmount
			chargeObject.BalanceTransaction = refundId
			//set refund object
			if refundAmount > 0 {
				//data
				datas := []RefundData{
					{
						Amount:             refundAmount,
						BalanceTransaction: refundId,
						Charge:             captureId,
						Created:            Timestamp(),
						Currency:           chargeObject.Currency,
						ID:                 "re_" + CreateChargeId(),
						Object:             "refund",
						Status:             "succeeded",
					},
				}
				//refund
				refunds := Refunds{
					TotalCount: 1,
					Object:     "list",
					HasMore:    false,
					URL:        "/v1/charges/" + captureId + "/refunds",
					Data:       datas,
				}
				//set refund
				chargeObject.Refunds = refunds
				chargeObject.Refunded = true
			}
			//success-write to stream
			json.NewEncoder(w).Encode(chargeObject)
			//success
			httpStatus = http.StatusOK
			//put object into cache
			cacheableObject := CacheObject{
				Status:      httpStatus,
				RequestId:   requestId,
				Charge:      chargeObject,
				Idempotency: idempotencyKey,
			}
			//cache item for next use
			captureCache.Set(captureId, cacheableObject, cache.DefaultExpiration)
			//print
			fmt.Println("CaptureHandler : mock object created")
		}
	} else {
		//end user is trying to access service with the same request format
		errorObjects := ErrorResponse{
			Error: ErrorObject{
				Type:    "invalid_request_error",
				Message: "No such charge: " + captureId,
				Param:   "id",
			},
		}
		//write error message
		fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObjects))
		//
		fmt.Println("CaptureHandler : No such charge: " + captureId)
	}
	//write http status
	w.WriteHeader(httpStatus)
}
