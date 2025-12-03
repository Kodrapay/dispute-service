package dto

type DisputeCreateRequest struct {
	TransactionReference string `json:"transaction_reference"`
	Reason               string `json:"reason"`
	Details              string `json:"details"`
}

type DisputeEvidenceRequest struct {
	EvidenceURL string `json:"evidence_url"`
	Note        string `json:"note"`
}

type DisputeResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Reference string `json:"reference"`
}
