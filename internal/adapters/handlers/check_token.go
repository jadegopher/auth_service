package handlers

import (
	"auth/internal/service"
	"auth/proto"
	"context"
)

// CheckToken TODO functinal tests
func (h *handlers) CheckToken(
	ctx context.Context,
	req *proto.CheckTokenRequest,
) (_ *proto.CheckTokenResponse, err error) {
	ok, err := h.authService.CheckToken(ctx, req.Token)
	if err != nil {
		if err == service.ErrAccountNotFound {
			return nil, err
		}
		return nil, ErrInternalServer
	}

	if !ok {

	}

	// TODO: refresh token
	return &proto.CheckTokenResponse{Token: req.Token}, nil
}
