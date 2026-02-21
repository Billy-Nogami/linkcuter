package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidCode       = errors.New("invalid code")
	ErrCodeAlreadyExists = errors.New("code already exists")
	ErrURLAlreadyExists  = errors.New("url already exists")
)
