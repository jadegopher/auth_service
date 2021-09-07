package handlers

import (
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
		return nil, ErrInternalServer
	}

	if !ok {
		return nil, ErrInvalidToken
	}

	// TODO: refresh token
	return &proto.CheckTokenResponse{Token: req.Token}, nil
}
