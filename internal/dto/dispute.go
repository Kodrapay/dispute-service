package dto

type DisputeCreateRequest struct {
	TransactionReference int `json:"transaction_reference"`
	Reason               string `json:"reason"`
	Details              string `json:"details"`
}

type DisputeEvidenceRequest struct {
	EvidenceURL int `json:"evidence_url"`
	Note        string `json:"note"`
}

type DisputeResponse struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Reference string `json:"reference"`
	MerchantID int    `json:"merchant_id,omitempty"`
	OpenedAt  string `json:"opened_at,omitempty"`
	ClosedAt  string `json:"closed_at,omitempty"`
}
