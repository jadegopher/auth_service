package handlers

import (
	"auth/internal/adapters"
	"auth/proto"
)

type handlers struct {
	authService adapters.AuthService
}

func New(authService adapters.AuthService) proto.AuthServiceServer {
	return &handlers{
		authService: authService,
	}
}
