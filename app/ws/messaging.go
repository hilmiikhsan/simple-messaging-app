package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/env"
)

func (s *service) ServeWSMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan models.MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true

		for {
			var message models.MessagePayload
			if err := c.ReadJSON(&message); err != nil {
				log.Error("Error reading message: ", err)
				break
			}

			message.Date = time.Now()

			err := s.messageRepository.InsertNewMessage(context.Background(), message)
			if err != nil {
				fmt.Println("Error inserting message: ", err)
			}

			broadcast <- message
		}
	}))

	go func() {
		for {
			message := <-broadcast
			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Error("Error sending message: ", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "4001"))))
}

func (s *service) GetMessageHistory(ctx context.Context) ([]models.MessagePayload, error) {
	resp, err := s.messageRepository.GetAllMessage(ctx)
	if err != nil {
		log.Error("Error getting all message: ", err)
		return resp, err
	}

	return resp, nil
}
