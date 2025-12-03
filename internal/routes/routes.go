package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/dispute-service/internal/handlers"
	"github.com/kodra-pay/dispute-service/internal/services"
)

func Register(app *fiber.App, service string) {
	health := handlers.NewHealthHandler(service)
	health.Register(app)

	svc := services.NewDisputeService()
	h := handlers.NewDisputeHandler(svc)
	api := app.Group("/disputes")
	api.Post("/", h.Create)
	api.Get("/:id", h.Get)
	api.Post("/:id/evidence", h.AddEvidence)
}
