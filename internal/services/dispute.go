package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kodra-pay/dispute-service/internal/dto"
	"github.com/kodra-pay/dispute-service/internal/models"
	"github.com/kodra-pay/dispute-service/internal/repositories"
)

type DisputeService struct {
	repo *repositories.DisputeRepository
}

func NewDisputeService(repo *repositories.DisputeRepository) *DisputeService { return &DisputeService{repo: repo} }

func (s *DisputeService) Create(ctx context.Context, req dto.DisputeCreateRequest) (dto.DisputeResponse, error) {
	if req.TransactionReference == "" {
		return dto.DisputeResponse{}, fmt.Errorf("transaction_reference is required")
	}
	dispute, err := s.repo.Create(ctx, req.TransactionReference, req.Reason)
	if err != nil {
		return dto.DisputeResponse{}, err
	}
	return dto.DisputeResponse{
		ID:         dispute.ID,
		Status:     string(dispute.Status),
		Reference:  dispute.Reference,
		MerchantID: dispute.MerchantID,
		OpenedAt:   dispute.OpenedAt.Format(time.RFC3339),
	}, nil
}

func (s *DisputeService) Get(ctx context.Context, id string) (dto.DisputeResponse, error) {
	dispute, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.DisputeResponse{}, err
	}
	if dispute == nil {
		return dto.DisputeResponse{}, nil
	}
	resp := dto.DisputeResponse{
		ID:         dispute.ID,
		Status:     string(dispute.Status),
		Reference:  dispute.Reference,
		MerchantID: dispute.MerchantID,
		OpenedAt:   dispute.OpenedAt.Format(time.RFC3339),
	}
	if dispute.ClosedAt != nil {
		resp.ClosedAt = dispute.ClosedAt.Format(time.RFC3339)
	}
	return resp, nil
}

func (s *DisputeService) AddEvidence(ctx context.Context, id string, req dto.DisputeEvidenceRequest) (dto.DisputeResponse, error) {
	if req.EvidenceURL == "" {
		return dto.DisputeResponse{}, fmt.Errorf("evidence_url is required")
	}
	dispute, err := s.repo.AddEvidence(ctx, id, models.Evidence{URL: req.EvidenceURL, Note: req.Note})
	if err != nil {
		return dto.DisputeResponse{}, err
	}
	if dispute == nil {
		return dto.DisputeResponse{}, nil
	}
	return dto.DisputeResponse{
		ID:         dispute.ID,
		Status:     string(dispute.Status),
		Reference:  dispute.Reference,
		MerchantID: dispute.MerchantID,
		OpenedAt:   dispute.OpenedAt.Format(time.RFC3339),
	}, nil
}

func (s *DisputeService) ListByMerchant(ctx context.Context, merchantID string, limit int) ([]dto.DisputeResponse, error) {
	list, err := s.repo.ListByMerchant(ctx, merchantID, limit)
	if err != nil {
		return nil, err
	}
	var resp []dto.DisputeResponse
	for _, d := range list {
		item := dto.DisputeResponse{
			ID:         d.ID,
			Status:     string(d.Status),
			Reference:  d.Reference,
			MerchantID: d.MerchantID,
			OpenedAt:   d.OpenedAt.Format(time.RFC3339),
		}
		if d.ClosedAt != nil {
			item.ClosedAt = d.ClosedAt.Format(time.RFC3339)
		}
		resp = append(resp, item)
	}
	return resp, nil
}
