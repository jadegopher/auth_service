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
	if err = h.authService.CheckToken(ctx, req.Token); err != nil {
		if err != service.ErrTokenExpired || err != service.ErrInvalidToken || err != service.ErrAccountNotFound {
			return nil, ErrInternalServer
		}
		return nil, err
	}

	// TODO: refresh token
	return &proto.CheckTokenResponse{Token: req.Token}, nil
}
