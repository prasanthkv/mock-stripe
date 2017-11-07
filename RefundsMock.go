package main

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

//
// https://stripe.com/docs/api#create_refund-charge
//
func RefundsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RefundsHandler : INIT")
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
		fmt.Println("RefundsHandler : idempotency found for key: " + idempotencyKey)
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
			fmt.Println("RefundsHandler : idempotency auth key found")
			//end user is trying to access capture with same idempotency as of auth.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges/" + captureId + "/refunds' vs '/v1/charges'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else if idempotency.Type == "capture" {
			fmt.Println("RefundsHandler : idempotency capture key found")
			//end user is trying to access capture with same idempotency as of capture.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges/" + captureId + "/refunds' vs '/v1/charges/" + captureId + "/capture'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else if idempotency.Type == "void" && idempotency.RequestHash != formHash {
			fmt.Println("RefundsHandler : idempotency void key found with md5:" + formHash)
			//end user is trying to access capture with same idempotency as of void with different form parameters.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used with the same parameters they were first used with. Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else {
			//valid request lets process
			fmt.Println("RefundsHandler : idempotency auth key cache")
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
		fmt.Println("RefundsHandler : new request with idempotency key: " + idempotencyKey)
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
	cachedObj, found := voidCache.Get(captureId)
	if found {
		fmt.Println("RefundsHandler : cache fault:" + idempotencyKey)
		//get capture object from cache
		cacheObject := cachedObj.(CacheObject)
		//copy original request id
		header.Set("original-request", cacheObject.RequestId)
		//write to stream
		if cacheObject.Status == 200 {
			json.NewEncoder(w).Encode(cacheObject.Refund)
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
	fmt.Println("RefundsHandler :First time request")
	//original request id and request id will be same this case
	header.Set("original-request", requestId)
	//evaluate auth & capture
	cachedObj, found = captureCache.Get(captureId)
	flow := "Auth"
	//
	if !found {
		flow = "Capture"
		//void auth before capture
		cachedObj, found = chargeCache.Get(captureId)
	}
	//
	// all set
	//
	if found {
		fmt.Println("RefundsHandler : " + flow + " found for " + captureId)
		//process
		chargeObject := (cachedObj.(CacheObject)).Charge
		//
		reqAmount, err := strconv.Atoi(FindFist(r.Form["amount"]))
		reqReason := FindFist(r.Form["reason"])
		//
		fmt.Println("RefundsHandler : Start", err)
		//
		if reqAmount <= 0 {
			//charge amount should not be less than requested amount
			errorObjects := ErrorResponse{
				Error: ErrorObject{
					Type:    "invalid_request_error",
					Message: "Amount must be at least 50 cents.",
					Param:   "amount",
				},
			}
			//write to stream
			json.NewEncoder(w).Encode(errorObjects)
		} else if chargeObject.Amount < reqAmount {
			//charge amount should be greater than requested amount
			errorObjects := ErrorResponse{
				Error: ErrorObject{
					Type:    "invalid_request_error",
					Message: "You cannot partially refund an uncaptured charge. Instead, capture the charge for an amount less than the original amount",
					Param:   "amount",
				},
			}
			//write to stream
			json.NewEncoder(w).Encode(errorObjects)
		} else if reqReason != "" && reqReason != "duplicate" && reqReason != "fraudulent" && reqReason != "requested_by_customer" {
			//reason is an enum
			errorObjects := ErrorResponse{
				Error: ErrorObject{
					Type:    "invalid_request_error",
					Message: "Invalid reason: must be one of duplicate, fraudulent, or requested_by_customer",
					Param:   "reason",
				},
			}
			//write to stream
			json.NewEncoder(w).Encode(errorObjects)
		} else {
			baakiAmount := chargeObject.Amount - reqAmount
			refundId := "txn_" + CreateChargeId()
			chargeObject.Captured = true
			//set refund object
			refundData := RefundData{
				ID:                 "re_" + CreateChargeId(),
				Object:             "refund",
				Amount:             reqAmount,
				BalanceTransaction: refundId,
				Charge:             captureId,
				Created:            Timestamp(),
				Currency:           chargeObject.Currency,
				Status:             "succeeded",
			}
			//success-write to stream
			json.NewEncoder(w).Encode(refundData)
			//status ok
			httpStatus = http.StatusOK
			//put object into cache
			cacheableObject := CacheObject{
				Status:      httpStatus,
				RequestId:   requestId,
				Refund:      refundData,
				Idempotency: idempotencyKey,
			}
			//cache item for next use
			voidCache.Set(captureId, cacheableObject, cache.DefaultExpiration)
			//print
			fmt.Println("RefundsHandler : mock object created with baaki" + strconv.Itoa(baakiAmount))
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
		fmt.Println("RefundsHandler : No such charge: " + captureId)
	}
	//write http status
	w.WriteHeader(httpStatus)
}
