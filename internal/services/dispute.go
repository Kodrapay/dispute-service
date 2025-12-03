package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/kodra-pay/dispute-service/internal/dto"
)

type DisputeService struct{}

func NewDisputeService() *DisputeService { return &DisputeService{} }

func (s *DisputeService) Create(_ context.Context, req dto.DisputeCreateRequest) dto.DisputeResponse {
	return dto.DisputeResponse{
		ID:        "dispute_" + uuid.NewString(),
		Status:    "open",
		Reference: req.TransactionReference,
	}
}

func (s *DisputeService) Get(_ context.Context, id string) dto.DisputeResponse {
	return dto.DisputeResponse{
		ID:     id,
		Status: "open",
	}
}

func (s *DisputeService) AddEvidence(_ context.Context, id string, req dto.DisputeEvidenceRequest) map[string]string {
	return map[string]string{"id": id, "evidence_url": req.EvidenceURL, "status": "received"}
}
