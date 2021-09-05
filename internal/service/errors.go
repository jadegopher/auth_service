package service

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrTokenExpired    = errors.New("authentication token expired")
	ErrInvalidToken    = errors.New("invalid authentication token")
)
