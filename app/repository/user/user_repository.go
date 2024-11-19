package user

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
)

func (r *repository) InsertNewUser(ctx context.Context, user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *repository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	var result models.User
	err := database.DB.Where("username = ?", user.Username).First(&result).Error
	if err != nil {
		log.Error("Failed to get user: ", err)
		return nil, err
	}

	return &result, nil
}
