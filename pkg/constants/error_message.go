package constants

import "errors"

var (
	ErrUsernameAlreadyExists       = errors.New("username already exists")
	ErrUsernameOrPasswordIncorrect = errors.New("username or password is incorrect")
)
