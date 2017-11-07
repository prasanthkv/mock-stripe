package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"

	"github.com/patrickmn/go-cache"
)

//
func AuthauthorizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthauthorizeHandler : INIT")
	//build the form
	r.ParseForm()
	fmt.Println("-----")
	fmt.Println(r.Form)
	fmt.Println("-----")
	httpStatus := http.StatusBadRequest
	//
	// Set all Headers
	//
	header := w.Header()
	requestId := CreateRequestId()
	header.Set("content-type", "application/json")
	header.Set("stripe-version", "mock-1.0")
	header.Set("request-id", requestId)
	//copy the idempotency key to response
	idempotencyKey := r.Header.Get("idempotency-key");
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
		fmt.Println("AuthauthorizeHandler : idempotency found for key: " + idempotencyKey)
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
		if idempotency.Type == "capture" {
			fmt.Println("AuthauthorizeHandler : idempotency capture key found")
			//end user is trying to access capture with same idempotency as of capture.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges' vs '/v1/charges/" + idempotency.ChargeId + "/capture'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else if idempotency.Type == "void" {
			fmt.Println("AuthauthorizeHandler : idempotency void key found")
			//end user is trying to access capture with same idempotency as of void with different form parameters.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges' vs '/v1/charges/" + idempotency.ChargeId + "/refunds'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else if idempotency.Type == "auth" && idempotency.RequestHash != formHash {
			fmt.Println("AuthauthorizeHandler : idempotency auth key found with md5:" + formHash)
			//end user is trying to access capture with same idempotency as of auth.
			errorObjects.Error.Message = "Keys for idempotent requests can only be used with the same parameters they were first used with. Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request."
		} else {
			fmt.Println("AuthauthorizeHandler : idempotency auth key cache")
			//this would be a request for cached auth
			exit = false
		}
		// idempotency error so exit
		if exit {
			fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObjects))
			//final http status code
			w.WriteHeader(httpStatus)
			//exit
			fmt.Println("AuthauthorizeHandler : IDP_EXIT")
			return
		}
	} else {
		//new request
		fmt.Println("AuthauthorizeHandler : new request with idempotency key: " + idempotencyKey)
		idempotency := Idempotency{
			Type:        "auth",
			RequestId:   requestId,
			RequestHash: formHash,
		}
		//set cache for next use
		idempotencyCache.Set(idempotencyKey, idempotency, cache.DefaultExpiration)
	}
	//
	//check for cached void object
	//
	chargeObj, found := authCache.Get(idempotencyKey)
	if found {
		fmt.Println("AuthauthorizeHandler : cache fault:" + idempotencyKey)
		//get capture object from cache
		cacheObject := chargeObj.(CacheObject)
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
		//print and exit
		fmt.Println("AuthauthorizeHandler : cache exit with http_status:" + strconv.Itoa(cacheObject.Status))
		return
	}
	//
	// First time request
	//
	fmt.Println("AuthauthorizeHandler :First time request")
	//original request id and request id will be same this case
	header.Set("original-request", requestId)
	//
	fmt.Println("AuthauthorizeHandler : Start")
	//evaluate auth & capture
	chargeRequest, errorObject, status := ValidateAndMapAuth(r)
	chargeObject := ChargeObject{}
	//declined
	if status != http.StatusOK {
		//write to stream
		json.NewEncoder(w).Encode(errorObject)
		//
		fmt.Println("AuthauthorizeHandler : exiting with validation error")
	} else {
		//new charge id
		newId := CreateChargeId()
		chargeId := "ch_" + newId
		cardId := "card_" + CreateCardId()
		//
		// create mock response objects
		//
		source := Source{
			ID:          cardId,
			Object:      "card",
			Fingerprint: CreateFingerPrint(),
			Funding:     "credit", //TODO Test card to funding type mapping
			Last4:       LastFour(chargeRequest.Source.Number),
			//
			Brand:   "Visa", //TODO support test card id to brand mapping
			Country: "US", // TODO support test card to country mapping
			//
			AddressCity:    chargeRequest.Source.AddressCity,
			AddressCountry: chargeRequest.Source.AddressCountry,
			AddressLine1:   chargeRequest.Source.AddressLine1,
			AddressLine2:   chargeRequest.Source.AddressLine2,
			AddressState:   chargeRequest.Source.AddressState,
			AddressZip:     chargeRequest.Source.AddressZip,
		}
		//
		//card to test flows https://stripe.com/docs/testing
		//

		//address line check
		if chargeRequest.Source.AddressLine1 != "" {
			if chargeRequest.Source.Number == 4000000000000028 {
				source.AddressLine1Check = "fail"
			} else if chargeRequest.Source.Number == 4000000000000044 {
				source.AddressLine1Check = "unavailable"
			} else {
				source.AddressLine1Check = "pass"
			}
		}
		//zip code check
		if chargeRequest.Source.AddressZip != "" {
			if chargeRequest.Source.Number == 4000000000000036 {
				source.AddressZipCheck = "fail"
			} else if chargeRequest.Source.Number == 4000000000000044 {
				source.AddressZipCheck = "unavailable"
			} else {
				source.AddressZipCheck = "pass"
			}
		}
		//If a CVC number is provided, the cvc_check fails.
		if chargeRequest.Source.CVV > 0 {
			if chargeRequest.Source.Number == 4000000000000101 {
				source.CvcCheck = "fail"
			} else {
				source.CvcCheck = "pass"
			}
		}

		//mock refunds
		refunds := Refunds{
			Object:     "list",
			HasMore:    false,
			TotalCount: 0,
			URL:        "/v1/charges/" + chargeId + "/refunds",
		}

		//create mock charge response
		chargeObject = ChargeObject{
			ID:             chargeId,
			Status:         "succeeded",
			Object:         "charge",
			Amount:         chargeRequest.Amount,
			AmountRefunded: 0,
			Captured:       false, //todo based on the capture request flag,
			Created:        Timestamp(),
			Currency:       chargeRequest.Currency,
			Description:    chargeRequest.Description,
			Livemode:       false,
			Paid:           true, // in case if there are no error
			Refunded:       false,
			Outcome:        successOutcome,
			Refunds:        refunds,
			Source:         source,
		}

		//Charge succeeds with a risk_level of elevated
		if chargeRequest.Source.Number == 4000000000000036 {
			chargeObject.Review = "prv_" + newId
			chargeObject.Outcome = elevatedOutcome
		}

		//return write to stream
		json.NewEncoder(w).Encode(chargeObject)
		//should be the last
		status = http.StatusOK
		//
		fmt.Println("AuthauthorizeHandler : mock object created")
	}
	//put object into cache
	cacheableObject := CacheObject{
		Type:        "auth",
		Status:      status,
		RequestId:   requestId,
		Charge:      chargeObject,
		Error:       errorObject,
	}
	//set item
	authCache.Set(idempotencyKey, cacheableObject, cache.DefaultExpiration)
	//should be the last
	w.WriteHeader(status)
	//the end
	fmt.Println(" AuthauthorizeHandler : END")
}
