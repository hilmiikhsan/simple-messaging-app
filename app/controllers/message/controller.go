package message

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
)

type service interface {
	GetMessageHistory(ctx context.Context) ([]models.MessagePayload, error)
}

type Controller struct {
	app     *fiber.App
	service service
}

func NewController(app *fiber.App, service service) *Controller {
	return &Controller{
		app:     app,
		service: service,
	}
}
