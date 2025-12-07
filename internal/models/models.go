package models

import "time"

type DisputeStatus string

const (
	StatusOpen        DisputeStatus = "open"
	StatusUnderReview DisputeStatus = "under_review"
	StatusResolved    DisputeStatus = "resolved"
	StatusRejected    DisputeStatus = "rejected"
)

type Evidence struct {
	URL  string `json:"url"`
	Note string `json:"note,omitempty"`
}

type Dispute struct {
	ID            int           `json:"id"`
	TransactionID int           `json:"transaction_id"`
	Status        DisputeStatus `json:"status"`
	Reason        string        `json:"reason,omitempty"`
	Evidence      []Evidence    `json:"evidence,omitempty"`
	OpenedAt      time.Time     `json:"opened_at"`
	ClosedAt      *time.Time    `json:"closed_at,omitempty"`
	Reference     string        `json:"reference,omitempty"`
	MerchantID    int           `json:"merchant_id,omitempty"`
}
