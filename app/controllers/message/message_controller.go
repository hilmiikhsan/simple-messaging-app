package message

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/response"
)

func (h *Controller) GetMessageHistory(ctx *fiber.Ctx) error {
	resp, err := h.service.GetMessageHistory(ctx.Context())
	if err != nil {
		errResponse := fmt.Errorf("failed to get message history: %v", err)
		log.Error(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, resp, fiber.StatusOK)
}
