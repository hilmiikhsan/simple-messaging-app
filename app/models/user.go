package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;type:varchar(20);not null;default:''" json:"username" validate:"required,min=6,max=20"`
	FullName  string    `gorm:"type:varchar(100);not null;default:''" json:"full_name" validate:"required,min=6,max=100"`
	Password  string    `gorm:"type:varchar(255);not null;default:''" json:"password,omitempty" validate:"required,min=6"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"-"`
}

func (l User) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginResponse struct {
	Username     string `json:"username" `
	FullName     string `json:"full_name" `
	Token        string `json:"token" `
	RefreshToken string `json:"refresh_token" `
}

type RefreshTokenResponse struct {
	Token string `json:"token" `
}
