package errors

import (
	"errors"
)

var Unauthorized error = errors.New("Unauthorized!")
var UserNotFound error = errors.New("User not found!")
var WrongPassword error = errors.New("Wrong password!")
var UserAlreadyActivated error = errors.New("User already activated!")
var UserAlreadyRegistered error = errors.New("User already registered!")
var UserNotActivated error = errors.New("User not activated!")
var PostNotFound error = errors.New("Post not found!")
