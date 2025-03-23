package storage

import (
	"errors"
)

var ErrorUserExists = errors.New("user already exists")
var ErrorUserNotFound = errors.New("user not found")
