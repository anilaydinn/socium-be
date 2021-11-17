package errors

import (
	"errors"
)

var UserNotFound error = errors.New("User not found!")
var WrongPasswordError = errors.New("Wrong password!")
