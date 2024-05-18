package shared

import "errors"

var (
	ErrNotExists          = errors.New("user does not exist")
	ErrUserWithEmailExist = errors.New("user with such email already exist")
	ErrDefaultInternal    = errors.New("there was problem with the request")
)
