package adapters

import "context"

type AuthService interface {
	RegisterGuestAccount(ctx context.Context, username string) (token string, err error)
	CheckToken(ctx context.Context, token string) (err error)
}
