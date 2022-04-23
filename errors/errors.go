package errors

import (
	"errors"
)

var Unauthorized error = errors.New("Unauthorized!")
var UserNotFound error = errors.New("User not found!")
var WrongPassword error = errors.New("Wrong password!")
