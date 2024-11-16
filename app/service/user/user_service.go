package user

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hilmiikhsan/simple-messaging-app/app/models"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/constants"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/jwt_token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) Register(ctx context.Context, user *models.User) (*models.User, error) {
	userData, err := s.userRepository.GetUser(ctx, user)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Failed to get user: ", err)
		return nil, err
	}

	if userData != nil {
		log.Error("username already exists")
		return nil, constants.ErrUsernameAlreadyExists
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password: ", err)
		return nil, err
	}
	user.Password = string(hashPassword)

	err = s.userRepository.InsertNewUser(ctx, user)
	if err != nil {
		log.Error("Failed to insert new user: ", err)
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *service) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	userData, err := s.userRepository.GetUser(ctx, &models.User{
		Username: req.Username,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Failed to get user: ", err)
		return nil, err
	}

	if userData == nil {
		log.Error("username or password is incorrect")
		return nil, constants.ErrUsernameOrPasswordIncorrect
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password))
	if err != nil {
		log.Error("failed to compare password: ", err)
		return nil, constants.ErrUsernameOrPasswordIncorrect
	}

	now := time.Now()

	token, err := jwt_token.GenerateToken(ctx, userData.Username, userData.FullName, constants.TokenType, now)
	if err != nil {
		log.Error("failed to generate token: ", err)
		return nil, err
	}

	refreshToken, err := jwt_token.GenerateToken(ctx, userData.Username, userData.FullName, constants.RefreshTokenType, now)
	if err != nil {
		log.Error("failed to generate refresh token: ", err)
		return nil, err
	}

	userSession := &models.UserSession{
		UserID:              userData.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(jwt_token.MapTypeToken[constants.TokenType]),
		RefreshTokenExpired: now.Add(jwt_token.MapTypeToken[constants.RefreshTokenType]),
	}

	err = s.userSessionRepository.InsertNewUserSession(ctx, userSession)
	if err != nil {
		log.Error("failed to insert new user session: ", err)
		return nil, err
	}

	return &models.LoginResponse{
		Username:     userData.Username,
		FullName:     userData.FullName,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	err := s.userSessionRepository.DeleteUserSessionByToken(ctx, token)
	if err != nil {
		log.Error("failed to delete user session: ", err)
		return err
	}

	return nil
}

func (s *service) RefreshToken(ctx context.Context, token string, tokenExpired time.Time, refreshToken string) error {
	err := s.userSessionRepository.UpdateUserSessionToken(ctx, token, tokenExpired, refreshToken)
	if err != nil {
		log.Error("failed to update user session token: ", err)
		return err
	}

	return nil
}
