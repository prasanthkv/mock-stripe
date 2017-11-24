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
    //evaluate card
    cardName := FindFist(r.Form["source[name]"])
    cardObject := FindFist(r.Form["source[object]"])
    cardCVV, err := strconv.Atoi(FindFist(r.Form["source[exp_year]"]))
    //
    cardNumber, err := strconv.ParseInt(FindFist(r.Form["source[number]"]), 10, 64)
    cardExpMonth, err := strconv.Atoi(FindFist(r.Form["source[exp_month]"]))
    cardExpYear, err := strconv.Atoi(FindFist(r.Form["source[exp_year]"]))
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
    if (cardExpMonth < 1 && cardExpMonth > 12) || (cardExpYear == time.Now().Year() && cardExpMonth < int(time.Now().Month())) {
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
        AddressCity:    FindFist(r.Form["source[address_city]"]),
        AddressCountry: FindFist(r.Form["source[address_country]"]),
        AddressLine1:   FindFist(r.Form["source[address_line1]"]),
        AddressLine2:   FindFist(r.Form["source[address_line2]"]),
        AddressState:   FindFist(r.Form["source[address_state]"]),
        AddressZip:     FindFist(r.Form["source[address_zip]"]),
    }
    //evaluate rest of the request
    reqCurrency := FindFist(r.Form["currency"])
    reqCapture, err := strconv.ParseBool(FindFist(r.Form["currency"]))
    reqAmount, err := strconv.Atoi(FindFist(r.Form["amount"]))
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
        Description: FindFist(r.Form["description"]),
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
    return chargeRequest, errorObjects, httpStatus
}

func FindFist(arrayValue []string) (string){
    if len(arrayValue) > 0{
        return arrayValue[0]
    }
    return  ""
}
