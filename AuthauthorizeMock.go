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

	cObject, found := authCache.Get(idempotencyKey)
	//check for idem request
	if found {
		header.Set("replay", "true")
		//
		cacheObject := cObject.(CacheObject)
		//
		header.Set("original-request", cacheObject.RequestId)
		//return
		if cacheObject.Status == 200{
			fmt.Fprintln(w, json.NewEncoder(w).Encode(cacheObject.Charge))
		}else{
			fmt.Fprintln(w, json.NewEncoder(w).Encode(cacheObject.Error))
		}
		//should be the last
		w.WriteHeader(http.StatusOK)
	} else {
		header.Set("original-request", requestId)
		chargeRequest, errorObject, status := ValidateAndMapAuth(r)
		chargeObject := ChargeObject{}
		//declined
		if status != http.StatusOK {
			fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObject))
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
				chargeObject.Review = "prv_"+newId
				chargeObject.Outcome = elevatedOutcome
			}

			//return
			fmt.Fprintln(w, json.NewEncoder(w).Encode(chargeObject))
			//should be the last
			status = http.StatusOK
			fmt.Println(" AuthauthorizeHandler : StatusOK")
		}
		//put object into cache
		cacheObject := CacheObject{
			RequestId: requestId,
			Charge:chargeObject,
			Error: errorObject,
			Status: status,
			Idempotency:idempotencyKey,
			RequestHash:"todo",
		}
		authCache.Set(idempotencyKey, cacheObject, cache.DefaultExpiration)
		//should be the last
		w.WriteHeader(status)
	}
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
