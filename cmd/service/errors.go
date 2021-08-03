package service

import "errors"

var (
	ErrInternalServer  = errors.New("internal server error")
	ErrAccountNotFound = errors.New("account not found")
	ErrTokenExpired    = errors.New("authentication token expired")
	ErrInvalidToken    = errors.New("invalid authentication token")
)
