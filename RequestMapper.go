package main

import (
	"net/http"
	"fmt"
	"strconv"
	"time"
)

//
func ValidateAndMapAuth(r *http.Request) (chargeRequest ChargeRequest, errorObject ErrorResponse, status int) {
	httpStatus := http.StatusPaymentRequired
	errorObjects := ErrorResponse{
		Error: ErrorObject{
		},
	}
	//build the form
	r.ParseForm()
	fmt.Println(r.Form)
	//evaluate card
	cardName := r.Form["source[name]"][0]
	cardObject := r.Form["source[object]"][0]
	cardCVV, err := strconv.Atoi(r.Form["source[exp_year]"][0])
	//
	cardNumber, err := strconv.ParseInt(r.Form["source[number]"][0], 10, 64)
	cardExpMonth, err := strconv.Atoi(r.Form["source[exp_month]"][0])
	cardExpYear, err := strconv.Atoi(r.Form["source[exp_year]"][0])
	//validate card number
	if cardNumber <= 0 {
		errorObjects.Error.Type = "invalid_request_error"
		errorObjects.Error.Message = "Could not find payment information."
		return chargeRequest, errorObjects, httpStatus
	} else if (cardNumber < 1000000000000000) && (cardNumber > 10000000000000000) {
		//card should be of 16 digit
		errorObjects.Error.Charge = "ch_" + CreateChargeId()
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "card_declined"
		errorObjects.Error.DeclineCode = "test_mode_live_card"
		errorObjects.Error.Message = "Your card was declined. Your request was in test mode, but used a non test card. For a list of valid test cards, visit: https://stripe.com/docs/testing."
		return chargeRequest, errorObjects, httpStatus
	}
	//evaluate exp month
	if ((cardExpMonth < 1 && cardExpMonth > 12) || (cardExpYear == time.Now().Year() && cardExpMonth < int(time.Now().Month()))) {
		//card should be of 16 digit
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Param = "exp_month"
		errorObjects.Error.Code = "invalid_expiry_month"
		errorObjects.Error.Message = "Your card's expiration month is invalid."
		return chargeRequest, errorObjects, httpStatus
	}
	//evaluate exp month
	if cardExpYear < time.Now().Year() {
		//card should be of 16 digit
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Param = "exp_month"
		errorObjects.Error.Code = "invalid_expiry_month"
		errorObjects.Error.Message = "Your card's expiration year is invalid."
		return chargeRequest, errorObjects, httpStatus
	}
	//log error
	print(err)
	//
	source := SourceObject{
		Name:           cardName,
		Number:         cardNumber,
		ExpMonth:       cardExpMonth,
		ExpYear:        cardExpYear,
		Object:         cardObject,
		CVV:            cardCVV,
		AddressCity:    r.Form["source[address_city]"][0],
		AddressCountry: r.Form["source[address_country]"][0],
		AddressLine1:   r.Form["source[address_line1]"][0],
		AddressLine2:   r.Form["source[address_line2]"][0],
		AddressState:   r.Form["source[address_state]"][0],
		AddressZip:     r.Form["source[address_zip]"][0],
	}
	//evaluate rest of the request
	reqCurrency := r.Form["currency"][0]
	reqCapture, err := strconv.ParseBool(r.Form["currency"][0])
	reqAmount, err := strconv.Atoi(r.Form["amount"][0])
	//amount should be more than 50c
	if reqAmount < 50 {
		errorObjects.Error.Type = "invalid_request_error"
		errorObjects.Error.Param = "amount"
		errorObjects.Error.Message = "Amount must be at least 50 cents"
		return chargeRequest, errorObjects, httpStatus
	}
	//ChargeRequest
	chargeRequest = ChargeRequest{
		Amount:      reqAmount,
		Capture:     reqCapture,
		Currency:    reqCurrency,
		Source:      source,
		Description: r.Form["description"][0],
	}
	//
	fmt.Println("Valid Request")
	//now look for card number
	switch cardNumber {
	case 4242424242424241:
		errorObjects.Error.Param = "number"
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "incorrect_number"
		errorObjects.Error.Message = "Your card number is incorrect."
	case 4000000000000119:
		errorObjects.Error.Param = "number"
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "processing_error"
		errorObjects.Error.Message = "An error occurred while processing your card. Try again in a little bit."
		errorObjects.Error.Charge = "ch_" + CreateChargeId()
	case 4000000000000069:
		errorObjects.Error.Param = "exp_month"
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "expired_card"
		errorObjects.Error.Message = "Your card has expired."
		errorObjects.Error.Charge = "ch_" + CreateChargeId()
	case 4000000000000127:
		errorObjects.Error.Param = "cvc"
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "incorrect_cvc"
		errorObjects.Error.Message = "Your card's security code is incorrect."
		errorObjects.Error.Charge = "ch_" + CreateChargeId()
	case 4100000000000019:
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "card_declined"
		errorObjects.Error.Message = "Your card was declined."
		errorObjects.Error.Charge = "ch_" + CreateChargeId()
		errorObjects.Error.DeclineCode = "fraudulent"
	case 4000000000000002:
		errorObjects.Error.Type = "card_error"
		errorObjects.Error.Code = "card_declined"
		errorObjects.Error.Message = "Your card was declined."
		errorObjects.Error.Charge = "ch_" + CreateChargeId()
		errorObjects.Error.DeclineCode = "generic_decline"
	default:
		httpStatus = http.StatusOK
	}
	return chargeRequest, errorObjects, httpStatus;
}
