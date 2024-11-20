package message

import (
	"fmt"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/response"
	"go.elastic.co/apm"
)

func (h *Controller) GetMessageHistory(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "GetMessageHistory", "controller")
	defer span.End()

	resp, err := h.service.GetMessageHistory(spanCtx)
	if err != nil {
		errResponse := fmt.Errorf("failed to get message history: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	return response.SendSuccessResponse(ctx, resp, fiber.StatusOK)
}
