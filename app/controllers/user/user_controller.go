package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/constants"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/jwt_token"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/response"
	"go.elastic.co/apm"
)

func (h *Controller) Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	req := new(models.User)

	err := ctx.BodyParser(req)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request body: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	err = req.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request body: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	resp, err := h.service.Register(spanCtx, req)
	if err != nil {
		errResponse := fmt.Errorf("failed to register new user: %v", err)
		log.Println(errResponse)

		if strings.Contains(err.Error(), constants.ErrUsernameAlreadyExists.Error()) {
			log.Println("username already exists")
			return response.SendFailureResponse(ctx, fiber.StatusConflict, err.Error(), nil)
		}

		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, resp, http.StatusCreated)
}

func (h *Controller) Login(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	req := new(models.LoginRequest)

	err := ctx.BodyParser(req)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request body: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	err = req.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request body: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	resp, err := h.service.Login(spanCtx, req)
	if err != nil {
		errResponse := fmt.Errorf("failed to login: %v", err)
		log.Println(errResponse)

		if strings.Contains(err.Error(), constants.ErrUsernameOrPasswordIncorrect.Error()) {
			log.Println("username or password is incorrect")
			return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
		}

		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, resp, http.StatusOK)
}

func (h *Controller) Logout(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Logout", "controller")
	defer span.End()

	token := ctx.Get("Authorization")

	err := h.service.Logout(spanCtx, token)
	if err != nil {
		errResponse := fmt.Errorf("failed to logout: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, nil, http.StatusOK)
}

func (h *Controller) RefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	now := time.Now()
	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwt_token.GenerateToken(ctx.Context(), username, fullName, constants.TokenType, now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = h.service.RefreshToken(spanCtx, token, now.Add(jwt_token.MapTypeToken[constants.TokenType]), refreshToken)
	if err != nil {
		errResponse := fmt.Errorf("failed to refresh token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := &models.RefreshTokenResponse{
		Token: token,
	}

	return response.SendSuccessResponse(ctx, resp, http.StatusOK)
}
