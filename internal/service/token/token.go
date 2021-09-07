package token

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
)

type tokenGen struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *tokenGen {
	return &tokenGen{logger: logger}
}

// Create creates signed jwt token
func (t *tokenGen) Create(payload map[string]interface{}, key interface{}) (_ string, err error) {
	// Initialize JWT token and set payload
	token := jwt.New()
	for k, v := range payload {
		if err = token.Set(k, v); err != nil {
			t.logger.Error("error token.Set", zap.Error(err), zap.String("key", k), zap.Any("value", v))
			return "", err
		}
	}

	// Sign token
	var signedToken []byte
	if signedToken, err = jwt.Sign(token, jwa.RS256, key); err != nil {
		t.logger.Error("error jwt.Sign", zap.Error(err))
		return "", err
	}

	return string(signedToken), nil
}
