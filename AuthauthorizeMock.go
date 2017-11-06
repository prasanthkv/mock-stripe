package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

//
func AuthauthorizeHandler(w http.ResponseWriter, r *http.Request) {
	//build the form
	r.ParseForm()
	fmt.Println(r.Form)
	//
	requestId := CreateRequestId()
	fmt.Println("AuthauthorizeHandler : Init")
	//all request are json
	header := w.Header()
	header.Set("request-id", requestId)
	header.Set("content-type", "application/json")
	header.Set("stripe-version", "mock-1.0")
	//copy the idempotency key to response
	idempotencyKey := r.Header.Get("idempotency-key")
	header.Set("idempotency-key", idempotencyKey)
	//make hash of form this will help maintain idempotency
	formHash := MD5Hash(r.Form)
	header.Set("request-md5", formHash)
	//
	cObject, found := authCache.Get(idempotencyKey)
	//check for idem request
	if found {
		cacheObject := cObject.(CacheObject)
		//
		if cacheObject.RequestHash == formHash && cacheObject.Type == "auth" {
			//
			header.Set("original-request", cacheObject.RequestId)
			//return /write to stream
			if cacheObject.Status == 200 {
				json.NewEncoder(w).Encode(cacheObject.Charge)
			} else {
				json.NewEncoder(w).Encode(cacheObject.Error)
			}
			//should be the last
			w.WriteHeader(cacheObject.Status)
		} else {
			//end user is trying to access service with the same request format
			errorObjects := ErrorResponse{
				Error: ErrorObject{
					Type : "idempotency_error",
					Message :"Keys for idempotent requests can only be used with the same parameters they were first used with. Try using a key other than 'key2' if you meant to execute a different request.",
				},
			}
			//write to stream
			json.NewEncoder(w).Encode(errorObjects)
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		header.Set("original-request", requestId)
		chargeRequest, errorObject, status := ValidateAndMapAuth(r)
		chargeObject := ChargeObject{}
		//declined
		if status != http.StatusOK {
			//write to stream
			json.NewEncoder(w).Encode(errorObject)
		} else {
			//new charge id
			newId := CreateChargeId()
			chargeId := "ch_" + newId
			cardId := "card_" + CreateCardId()

			source := Source{
				ID:          cardId,
				Object:      "card",
				Fingerprint: CreateFingerPrint(),
				Funding:     "credit",
				Last4:       LastFour(chargeRequest.Source.Number),
				//
				Brand:   "Visa",
				Country: "US",
				//
				AddressCity:    chargeRequest.Source.AddressCity,
				AddressCountry: chargeRequest.Source.AddressCountry,
				AddressLine1:   chargeRequest.Source.AddressLine1,
				AddressLine2:   chargeRequest.Source.AddressLine2,
				AddressState:   chargeRequest.Source.AddressState,
				AddressZip:     chargeRequest.Source.AddressZip,
			}
			//set pass
			if chargeRequest.Source.AddressLine1 != "" {
				if chargeRequest.Source.Number == 4000000000000028 {
					source.AddressLine1Check = "fail"
				} else if chargeRequest.Source.Number == 4000000000000044 {
					source.AddressLine1Check = "unavailable"
				} else {
					source.AddressLine1Check = "pass"
				}
			}
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

			refunds := Refunds{
				Object:     "list",
				HasMore:    false,
				TotalCount: 0,
				URL:        "/v1/charges/" + chargeId + "/refunds",
			}

			//create response
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
			//cache charge
			chargeCache.Set(chargeId, chargeObject, cache.DefaultExpiration)
		}
		//put object into cache
		cacheObject := CacheObject{
			RequestId:   requestId,
			Charge:      chargeObject,
			Error:       errorObject,
			Status:      status,
			Idempotency: idempotencyKey,
			RequestHash: formHash,
			Type: "auth",
		}
		//set item
		authCache.Set(idempotencyKey, cacheObject, cache.DefaultExpiration)
		//should be the last
		w.WriteHeader(status)
	}
	//the end
	fmt.Println(" AuthauthorizeHandler : END ")
}

func LastFour(value int64) (trimValue string) {
	// Take substring of first word with runes.
	// ... This handles any kind of rune in the string.
	tempString := strconv.FormatInt(value, 10)
	length := len(tempString)
	//
	if length > 4 {
		runes := []rune(tempString)
		safeSubstring := string(runes[(length - 4):length])
		return safeSubstring;
	}
	return tempString
}

func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
