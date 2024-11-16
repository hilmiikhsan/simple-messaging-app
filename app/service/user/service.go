package user

import (
	"context"
	"time"

	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
)

type userRepository interface {
	InsertNewUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
}

type UserSessionRepository interface {
	InsertNewUserSession(ctx context.Context, session *models.UserSession) error
	GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error)
	DeleteUserSessionByToken(ctx context.Context, token string) error
	UpdateUserSessionToken(ctx context.Context, token string, tokenExpired time.Time, refreshToken string) error
}

type service struct {
	cfg                   *database.Config
	userRepository        userRepository
	userSessionRepository UserSessionRepository
}

func NewService(cfg *database.Config, userRepository userRepository, userSessionRepository UserSessionRepository) *service {
	return &service{
		cfg:                   cfg,
		userRepository:        userRepository,
		userSessionRepository: userSessionRepository,
	}
}
