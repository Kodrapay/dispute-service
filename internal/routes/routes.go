package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/dispute-service/internal/handlers"
	"github.com/kodra-pay/dispute-service/internal/config"
	"github.com/kodra-pay/dispute-service/internal/repositories"
	"github.com/kodra-pay/dispute-service/internal/services"
)

func Register(app *fiber.App, serviceName string, cfg config.Config) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	repo, err := repositories.NewDisputeRepository(cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}
	svc := services.NewDisputeService(repo)
	handler := handlers.NewDisputeHandler(svc)

	app.Post("/disputes", handler.Create)
	app.Get("/disputes/:id", handler.Get)
	app.Post("/disputes/:id/evidence", handler.AddEvidence)
	app.Get("/disputes", handler.ListByMerchant)
}
