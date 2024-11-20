package message

import (
	"context"

	"log"

	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"go.elastic.co/apm"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) InsertNewMessage(ctx context.Context, data models.MessagePayload) error {
	span, _ := apm.StartSpan(ctx, "InsertNewMessage", "repository")
	defer span.End()

	_, err := database.MongoDB.InsertOne(ctx, data)
	return err
}

func (r *repository) GetAllMessage(ctx context.Context) ([]models.MessagePayload, error) {
	span, _ := apm.StartSpan(ctx, "GetAllMessage", "repository")
	defer span.End()

	var (
		err  error
		resp []models.MessagePayload
	)

	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error getting all message: ", err)
		return resp, err
	}

	for cursor.Next(ctx) {
		payload := models.MessagePayload{}

		if err := cursor.Decode(&payload); err != nil {
			log.Println("Error decoding message: ", err)
			return resp, err
		}

		resp = append(resp, payload)
	}

	return resp, nil
}
