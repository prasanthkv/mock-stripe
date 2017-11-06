package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"

	"strconv"
)

func CaptureHandler(w http.ResponseWriter, r *http.Request) {
	//build the form
	r.ParseForm()
	fmt.Println(r.Form)
	httpStatus := http.StatusBadRequest
	//all request are json
	header := w.Header()
	requestId := CreateRequestId()
	header.Set("request-id", requestId)
	header.Set("content-type", "application/json")
	header.Set("stripe-version", "mock-1.0")
	//capture id
	vars := mux.Vars(r)
	captureId := vars["id"]
	//copy the idempotency key to response
	idempotencyKey := r.Header.Get("idempotency-key");
	header.Set("idempotency-key", idempotencyKey)
	header.Set("original-capture-id", captureId)
	//evaluate idempotency key
	//make hash of form this will help maintain idempotency
	formHash := MD5Hash(r.Form)
	header.Set("request-md5", formHash)
	_, found := authCache.Get(idempotencyKey)
	//check for auth idempotency_key
	if found {
		//the end
		fmt.Println("auth key found")
		//end user is trying to access capture with same idempotency as of auth.
		errorObjects := ErrorResponse{
			Error: ErrorObject{
				Type:    "idempotency_error",
				Message: "Keys for idempotent requests can only be used for the same endpoint they were first used for ('/v1/charges/" + captureId + "/capture' vs '/v1/charges'). Try using a key other than '" + idempotencyKey + "' if you meant to execute a different request.",
			},
		}
		fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObjects))
		//final http status code
		w.WriteHeader(httpStatus)
		return;
	}
	//process capture
	fmt.Println("RequestHash :new")
	//evaluate charge
	chargeObj, found := chargeCache.Get(captureId)
	//
	if found {
		//process
		chargeObject := chargeObj.(ChargeObject)
		reqAmount, err := strconv.Atoi(FindFist(r.Form["amount"]))
		//
		print(err)
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
					RefundData{
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
			}
			//success-write to stream
			json.NewEncoder(w).Encode(chargeObject)
			httpStatus = http.StatusOK
		}
	} else {
		//end user is trying to access service with the same request format
		errorObjects := ErrorResponse{
			Error: ErrorObject{
				Type:    "invalid_request_error",
				Message: "No such charge: undefined",
				Param:   "id",
			},
		}
		fmt.Fprintln(w, json.NewEncoder(w).Encode(errorObjects))
	}
	//return
	//should be the last
	w.WriteHeader(httpStatus)
}
