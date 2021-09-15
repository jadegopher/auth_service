package handlers

import (
	"auth/internal/adapters"
)

type handlers struct {
	authService adapters.AuthService
}

func New(authService adapters.AuthService) *handlers {
	return &handlers{
		authService: authService,
	}
}
