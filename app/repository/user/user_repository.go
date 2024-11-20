package user

import (
	"context"

	"log"

	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"go.elastic.co/apm"
)

func (r *repository) InsertNewUser(ctx context.Context, user *models.User) error {
	span, _ := apm.StartSpan(ctx, "InsertNewUser", "repository")
	defer span.End()

	return database.DB.Create(user).Error
}

func (r *repository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	span, _ := apm.StartSpan(ctx, "GetUser", "repository")
	defer span.End()

	var result models.User
	err := database.DB.Where("username = ?", user.Username).First(&result).Error
	if err != nil {
		log.Println("Failed to get user: ", err)
		return nil, err
	}

	return &result, nil
}
