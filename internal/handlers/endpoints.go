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
	resp, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *DisputeHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	resp, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if resp.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "dispute not found")
	}
	return c.JSON(resp)
}

func (h *DisputeHandler) AddEvidence(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.DisputeEvidenceRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp, err := h.svc.AddEvidence(c.Context(), id, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if resp.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "dispute not found")
	}
	return c.JSON(resp)
}

func (h *DisputeHandler) ListByMerchant(c *fiber.Ctx) error {
	merchantID := c.Query("merchant_id")
	if merchantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "merchant_id is required")
	}
	limit := c.QueryInt("limit", 50)
	resp, err := h.svc.ListByMerchant(c.Context(), merchantID, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(resp)
}
