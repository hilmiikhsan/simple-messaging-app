package message

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) InsertNewMessage(ctx context.Context, data models.MessagePayload) error {
	_, err := database.MongoDB.InsertOne(ctx, data)
	return err
}

func (r *repository) GetAllMessage(ctx context.Context) ([]models.MessagePayload, error) {
	var (
		err  error
		resp []models.MessagePayload
	)

	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		log.Error("Error getting all message: ", err)
		return resp, err
	}

	for cursor.Next(ctx) {
		payload := models.MessagePayload{}

		if err := cursor.Decode(&payload); err != nil {
			log.Error("Error decoding message: ", err)
			return resp, err
		}

		resp = append(resp, payload)
	}

	return resp, nil
}
