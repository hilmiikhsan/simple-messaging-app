package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/app/service/user"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/jwt_token"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/response"
)

type Middleware struct {
	userSessionRepository user.UserSessionRepository
}

func NewMiddleware(userSessionRepository user.UserSessionRepository) *Middleware {
	return &Middleware{
		userSessionRepository: userSessionRepository,
	}
}

func (m *Middleware) MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")

	if auth == "" {
		log.Error("Authorization header is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	_, err := m.userSessionRepository.GetUserSessionByToken(ctx.Context(), auth)
	if err != nil {
		log.Error("Failed to get user session by token: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Error("Failed to validate token: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Error("JWT token is expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

func (m *Middleware) MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")

	if auth == "" {
		log.Error("Authorization header is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Error("Failed to validate token: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Error("JWT token is expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}
