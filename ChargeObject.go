package main

type ChargeObject struct {
    ID                  string      `json:"id,omitempty"`
    Object              string      `json:"object,omitempty"`
    Created             int64       `json:"created,omitempty"`
    Amount              int         `json:"amount,omitempty"`
    Livemode            bool        `json:"livemode"`
    Refunded            bool        `json:"refunded"`
    Captured            bool        `json:"captured"`
    Paid                bool        `json:"paid"`
    Status              string      `json:"status,omitempty"`
    AmountRefunded      int         `json:"amount_refunded,omitempty"`
    Application         interface{} `json:"application,omitempty"`
    ApplicationFee      interface{} `json:"application_fee,omitempty"`
    BalanceTransaction  interface{} `json:"balance_transaction,omitempty"`
    Currency            string      `json:"currency,omitempty"`
    Customer            interface{} `json:"customer,omitempty"`
    Description         string      `json:"description,omitempty"`
    Destination         interface{} `json:"destination,omitempty"`
    Dispute             interface{} `json:"dispute,omitempty"`
    FailureCode         interface{} `json:"failure_code,omitempty"`
    FailureMessage      interface{} `json:"failure_message,omitempty"`
    Invoice             interface{} `json:"invoice,omitempty"`
    OnBehalfOf          interface{} `json:"on_behalf_of,omitempty"`
    Order               interface{} `json:"order,omitempty"`
    ReceiptEmail        interface{} `json:"receipt_email,omitempty"`
    ReceiptNumber       interface{} `json:"receipt_number,omitempty"`
    Review              interface{} `json:"review,omitempty"`
    Shipping            interface{} `json:"shipping,omitempty"`
    SourceTransfer      interface{} `json:"source_transfer,omitempty"`
    StatementDescriptor interface{} `json:"statement_descriptor,omitempty"`
    TransferGroup       interface{} `json:"transfer_group,omitempty"`
    FraudDetails struct {
    } `json:"fraud_details,omitempty"`
    Outcome Outcome `json:"outcome,omitempty"`
    Refunds Refunds `json:"refunds,omitempty"`
    Source  Source  `json:"source,omitempty"`
    // Metadata map `json:"metadata,omitempty"`
}

type FraudDetails struct {
}

type Outcome struct {
    NetworkStatus string `json:"network_status,omitempty"`
    Reason        string `json:"reason,omitempty"`
    RiskLevel     string `json:"risk_level,omitempty"`
    SellerMessage string `json:"seller_message,omitempty"`
    Type          string `json:"type,omitempty"`
}

type Refunds struct {
    Object     string       `json:"object,omitempty"`
    HasMore    bool         `json:"has_more"`
    TotalCount int          `json:"total_count"`
    URL        string       `json:"url,omitempty"`
    Data       []RefundData `json:"data,omitempty"`
}

type Source struct {
    ID                string `json:"id,omitempty"`
    Object            string `json:"object,omitempty"`
    AddressCity       string `json:"address_city,omitempty"`
    AddressCountry    string `json:"address_country,omitempty"`
    AddressLine1      string `json:"address_line1,omitempty"`
    AddressLine2      string `json:"address_line2,omitempty"`
    AddressState      string `json:"address_state,omitempty"`
    AddressZip        string `json:"address_zip,omitempty"`
    AddressLine1Check string `json:"address_line1_check,omitempty"`
    AddressZipCheck   string `json:"address_zip_check,omitempty"`
    CvcCheck          string `json:"cvc_check,omitempty"`
    Brand             string `json:"brand,omitempty"`
    Country           string `json:"country,omitempty"`
    Customer          string `json:"customer,omitempty"`
    DynamicLast4      string `json:"dynamic_last4,omitempty"`
    ExpMonth          int    `json:"exp_month,omitempty"`
    ExpYear           int    `json:"exp_year,omitempty"`
    Fingerprint       string `json:"fingerprint,omitempty"`
    Funding           string `json:"funding,omitempty"`
    Last4             string `json:"last4,omitempty"`
    Metadata struct {
        Key1 string `json:"key_1,omitempty"`
        Key2 string `json:"key_2,omitempty"`
    } `json:"metadata,omitempty"`
    Name string `json:"name,omitempty"`
}

type ErrorResponse struct {
    Error ErrorObject `json:"error"`
}

type ErrorObject struct {
    Type        string `json:"type,omitempty"`
    Message     string `json:"message,omitempty"`
    Param       string `json:"param,omitempty"`
    Code        string `json:"code,omitempty"`
    Charge      string `json:"charge,omitempty"`
    DeclineCode string `json:"decline_code,omitempty"`
}

type RefundData struct {
    Amount             int    `json:"amount,omitempty"`
    BalanceTransaction string `json:"balance_transaction,omitempty"`
    Charge             string `json:"charge,omitempty"`
    Created            int64  `json:"created,omitempty"`
    Currency           string `json:"currency,omitempty"`
    ID                 string `json:"id,omitempty"`
    Object             string `json:"object,omitempty"`
    Status             string `json:"status,omitempty"`
    Reason             string `json:"reason,omitempty"`
}
