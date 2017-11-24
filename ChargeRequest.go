package main

type ChargeRequest struct {
    Amount      int
    Currency    string
    Capture     bool
    Source      SourceObject
    Description string
}

type SourceObject struct {
    Object string
    //card info
    Name     string
    Number   int64
    ExpMonth int
    ExpYear  int
    CVV      int
    //address
    AddressCity    string
    AddressCountry string
    AddressLine1   string
    AddressLine2   string
    AddressState   string
    AddressZip     string
}
