package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/dispute-service/internal/dto"
	"github.com/kodra-pay/dispute-service/internal/services"
)

type DisputeHandler struct {
	svc *services.DisputeService
}

func NewDisputeHandler(svc *services.DisputeService) *DisputeHandler {
	return &DisputeHandler{svc: svc}
}

func (h *DisputeHandler) Create(c *fiber.Ctx) error {
	var req dto.DisputeCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return c.JSON(h.svc.Create(c.Context(), req))
}

func (h *DisputeHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(h.svc.Get(c.Context(), id))
}

func (h *DisputeHandler) AddEvidence(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.DisputeEvidenceRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return c.JSON(h.svc.AddEvidence(c.Context(), id, req))
}
