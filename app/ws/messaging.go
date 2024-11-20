package ws

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/env"
	"go.elastic.co/apm"
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
				log.Println("Error reading message: ", err)
				break
			}

			tx := apm.DefaultTracer.StartTransaction("Send Message", "ws")
			ctx := apm.ContextWithTransaction(context.Background(), tx)

			message.Date = time.Now()

			err := s.messageRepository.InsertNewMessage(ctx, message)
			if err != nil {
				fmt.Println("Error inserting message: ", err)
			}

			tx.End()

			broadcast <- message
		}
	}))

	go func() {
		for {
			message := <-broadcast
			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Println("Error sending message: ", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "4001"))))
}

func (s *service) GetMessageHistory(ctx context.Context) ([]models.MessagePayload, error) {
	span, _ := apm.StartSpan(ctx, "GetMessageHistory", "service")
	defer span.End()

	resp, err := s.messageRepository.GetAllMessage(ctx)
	if err != nil {
		log.Println("Error getting all message: ", err)
		return resp, err
	}

	return resp, nil
}
