package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type UserSession struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement"`
	UserID              uint      `gorm:"type:int;not null;default:0" json:"user_id"`
	Token               string    `gorm:"type:varchar(255);not null;default:''" json:"token"`
	RefreshToken        string    `gorm:"type:varchar(255);not null;default:''" json:"refresh_token"`
	TokenExpired        time.Time `gorm:"type:datetime;not null" json:"token_expired"`
	RefreshTokenExpired time.Time `gorm:"type:datetime;not null" json:"refresh_token_expired"`
	CreatedAt           time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt           time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"-"`
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
