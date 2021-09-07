package handlers

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidToken   = errors.New("invalid or expired token")
)
