package domain

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrDuplicateData = errors.New("resource already exists")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrForbidden     = errors.New("not permitted to perform this action")
)