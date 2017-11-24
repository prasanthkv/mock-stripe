package main

var successOutcome = Outcome{
    NetworkStatus: "approved_by_network",
    RiskLevel:     "normal",
    SellerMessage: "Payment complete.",
    Type:          "authorized",
}

var elevatedOutcome = Outcome{
    NetworkStatus: "approved_by_network",
    Reason:        "elevated_risk_level",
    RiskLevel:     "elevated",
    SellerMessage: "Stripe evaluated this charge as having elevated risk, and placed it in your manual review queue.",
    Type:          "manual_review",
}
