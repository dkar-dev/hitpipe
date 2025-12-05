package user

import "errors"

var ErrUserNotFound = errors.New("user not found")

var ErrTokenAlreadyExists = errors.New("token already exists")
