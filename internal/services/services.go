package services

import "github.com/kodra-pay/dispute-service/internal/repositories"

type Services struct {
	Disputes *DisputeService
}

func New(repo *repositories.DisputeRepository) *Services {
	return &Services{
		Disputes: NewDisputeService(repo),
	}
}
