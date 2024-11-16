package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/constants"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/jwt_token"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/response"
)

func (h *Controller) Register(ctx *fiber.Ctx) error {
	req := new(models.User)

	err := ctx.BodyParser(req)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request body: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	err = req.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request body: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	resp, err := h.service.Register(ctx.Context(), req)
	if err != nil {
		errResponse := fmt.Errorf("failed to register new user: %v", err)
		log.Error(errResponse)

		if strings.Contains(err.Error(), constants.ErrUsernameAlreadyExists.Error()) {
			log.Error("username already exists")
			return response.SendFailureResponse(ctx, fiber.StatusConflict, err.Error(), nil)
		}

		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, resp, http.StatusCreated)
}

func (h *Controller) Login(ctx *fiber.Ctx) error {
	req := new(models.LoginRequest)

	err := ctx.BodyParser(req)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request body: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	err = req.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request body: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, errResponse.Error(), nil)
	}

	resp, err := h.service.Login(ctx.Context(), req)
	if err != nil {
		errResponse := fmt.Errorf("failed to login: %v", err)
		log.Error(errResponse)

		if strings.Contains(err.Error(), constants.ErrUsernameOrPasswordIncorrect.Error()) {
			log.Error("username or password is incorrect")
			return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
		}

		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, resp, http.StatusOK)
}

func (h *Controller) Logout(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	err := h.service.Logout(ctx.Context(), token)
	if err != nil {
		errResponse := fmt.Errorf("failed to logout: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, nil, http.StatusOK)
}

func (h *Controller) RefreshToken(ctx *fiber.Ctx) error {
	now := time.Now()
	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwt_token.GenerateToken(ctx.Context(), username, fullName, constants.TokenType, now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = h.service.RefreshToken(ctx.Context(), token, now.Add(jwt_token.MapTypeToken[constants.TokenType]), refreshToken)
	if err != nil {
		errResponse := fmt.Errorf("failed to refresh token: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := &models.RefreshTokenResponse{
		Token: token,
	}

	return response.SendSuccessResponse(ctx, resp, http.StatusOK)
}
