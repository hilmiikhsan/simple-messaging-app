package user_session

import (
	"context"
	"time"

	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"go.elastic.co/apm"
)

func (r *repository) InsertNewUserSession(ctx context.Context, session *models.UserSession) error {
	span, _ := apm.StartSpan(ctx, "InsertNewUserSession", "repository")
	defer span.End()

	return database.DB.Create(session).Error
}

func (r *repository) GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error) {
	span, _ := apm.StartSpan(ctx, "GetUserSessionByToken", "repository")
	defer span.End()

	var (
		resp models.UserSession
		err  error
	)

	err = database.DB.Where("token = ?", token).Last(&resp).Error
	return resp, err
}

func (r *repository) DeleteUserSessionByToken(ctx context.Context, token string) error {
	span, _ := apm.StartSpan(ctx, "DeleteUserSessionByToken", "repository")
	defer span.End()

	return database.DB.Exec(queryDeleteUserSessionByToken, token).Error
}

func (r *repository) UpdateUserSessionToken(ctx context.Context, token string, tokenExpired time.Time, refreshToken string) error {
	span, _ := apm.StartSpan(ctx, "UpdateUserSessionToken", "repository")
	defer span.End()

	return database.DB.Exec(queryUpdteUserSessionToken, token, tokenExpired, refreshToken).Error
}
