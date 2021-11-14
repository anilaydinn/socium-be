package errors

import (
	"errors"
)

var UserNotFound error = errors.New("User not found!")
