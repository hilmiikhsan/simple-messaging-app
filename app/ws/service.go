package ws

import (
	"context"

	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
)

type messageRepository interface {
	InsertNewMessage(ctx context.Context, data models.MessagePayload) error
	GetAllMessage(ctx context.Context) ([]models.MessagePayload, error)
}

type service struct {
	cfg               *database.Config
	messageRepository messageRepository
}

func NewService(cfg *database.Config, messageRepository messageRepository) *service {
	return &service{
		cfg:               cfg,
		messageRepository: messageRepository,
	}
}
