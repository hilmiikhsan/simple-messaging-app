package user

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
)

type service interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string, tokenExpired time.Time, refreshToken string) error
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
