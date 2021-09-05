package handlers

import (
	"auth/proto"
	"context"
)

// RegisterGuestAccount TODO functional test
func (h *handlers) RegisterGuestAccount(
	ctx context.Context,
	req *proto.RegisterGuestAccountRequest,
) (_ *proto.RegisterGuestAccountResponse, err error) {
	token, err := h.authService.RegisterGuestAccount(ctx, req.Name)
	if err != nil {
		return nil, ErrInternalServer
	}

	return &proto.RegisterGuestAccountResponse{Token: token}, nil
}
