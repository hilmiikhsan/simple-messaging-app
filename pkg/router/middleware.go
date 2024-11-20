package router

import (
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/service/user"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/jwt_token"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/response"
	"go.elastic.co/apm"
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
	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareValidateAuth", "middleware")
	defer span.End()

	auth := ctx.Get("Authorization")

	if auth == "" {
		log.Println("Authorization header is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	_, err := m.userSessionRepository.GetUserSessionByToken(spanCtx, auth)
	if err != nil {
		log.Println("Failed to get user session by token: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(spanCtx, auth)
	if err != nil {
		log.Println("Failed to validate token: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("JWT token is expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

func (m *Middleware) MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareRefreshToken", "middleware")
	defer span.End()

	auth := ctx.Get("Authorization")

	if auth == "" {
		log.Println("Authorization header is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(spanCtx, auth)
	if err != nil {
		log.Println("Failed to validate token: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("JWT token is expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}
