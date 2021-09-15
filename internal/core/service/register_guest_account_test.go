package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"auth/mocks"
)

type testMockKey struct{}

func (t *testMockKey) String() string {
	return "key"
}

func (t *testMockKey) Interface() interface{} {
	return "key"
}

func Test_service_RegisterGuestAccount(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name string
		args args
		mock func(
			users *mocks.MockIUsers,
			sessions *mocks.MockISessions,
			keys *mocks.MockKeyGenerator,
			token *mocks.MockTokenCreator,
		)
		want    string
		wantErr bool
	}{
		{
			name: "success: everything is allright",
			args: args{
				ctx:      context.Background(),
				username: "somebody was told me",
			},
			mock: func(
				users *mocks.MockIUsers,
				sessions *mocks.MockISessions,
				keys *mocks.MockKeyGenerator,
				token *mocks.MockTokenCreator,
			) {
				mockKey := &testMockKey{}

				keys.EXPECT().Generate().Return(mockKey, nil)

				users.EXPECT().Insert(
					context.Background(),
					"somebody was told me",
					mockKey.String(),
				).Return(int64(1), nil)

				token.EXPECT().Create(
					map[string]interface{}{
						issuedAtField: time.Now().UTC().Unix(),
						userIDField:   int64(1),
					},
					mockKey.Interface(),
				).Return("signed key", nil)

				sessions.EXPECT().Insert("signed key", int64(1), tokenLifeTime).Return(nil)
			},
			want:    "signed key",
			wantErr: false,
		},
		{
			name: "failed: key generator error",
			args: args{
				ctx:      context.Background(),
				username: "somebody was told me",
			},
			mock: func(
				users *mocks.MockIUsers,
				sessions *mocks.MockISessions,
				keys *mocks.MockKeyGenerator,
				token *mocks.MockTokenCreator,
			) {
				keys.EXPECT().Generate().Return(nil, errors.New("some error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed: user info insert error",
			args: args{
				ctx:      context.Background(),
				username: "somebody was told me",
			},
			mock: func(
				users *mocks.MockIUsers,
				sessions *mocks.MockISessions,
				keys *mocks.MockKeyGenerator,
				token *mocks.MockTokenCreator,
			) {
				mockKey := &testMockKey{}
				keys.EXPECT().Generate().Return(mockKey, nil)

				users.EXPECT().Insert(
					context.Background(),
					"somebody was told me",
					mockKey.String(),
				).Return(int64(0), errors.New("some error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed: create token error",
			args: args{
				ctx:      context.Background(),
				username: "somebody was told me",
			},
			mock: func(
				users *mocks.MockIUsers,
				sessions *mocks.MockISessions,
				keys *mocks.MockKeyGenerator,
				token *mocks.MockTokenCreator,
			) {
				mockKey := &testMockKey{}
				keys.EXPECT().Generate().Return(mockKey, nil)

				users.EXPECT().Insert(
					context.Background(),
					"somebody was told me",
					mockKey.String(),
				).Return(int64(1), nil)

				token.EXPECT().Create(
					map[string]interface{}{
						issuedAtField: time.Now().UTC().Unix(),
						userIDField:   int64(1),
					},
					mockKey.Interface(),
				).Return("", errors.New("some error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed: insert token error",
			args: args{
				ctx:      context.Background(),
				username: "somebody was told me",
			},
			mock: func(
				users *mocks.MockIUsers,
				sessions *mocks.MockISessions,
				keys *mocks.MockKeyGenerator,
				token *mocks.MockTokenCreator,
			) {
				mockKey := &testMockKey{}

				keys.EXPECT().Generate().Return(mockKey, nil)

				users.EXPECT().Insert(
					context.Background(),
					"somebody was told me",
					mockKey.String(),
				).Return(int64(1), nil)

				token.EXPECT().Create(
					map[string]interface{}{
						issuedAtField: time.Now().UTC().Unix(),
						userIDField:   int64(1),
					},
					mockKey.Interface(),
				).Return("signed key", nil)

				sessions.EXPECT().Insert(
					"signed key",
					int64(1),
					tokenLifeTime,
				).Return(errors.New("some error"))
			},
			want:    "",
			wantErr: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersDB := mocks.NewMockIUsers(ctrl)
	sessionsDB := mocks.NewMockISessions(ctrl)
	tokenCreator := mocks.NewMockTokenCreator(ctrl)
	keyGenerator := mocks.NewMockKeyGenerator(ctrl)

	s := &Service{
		logger:       zap.New(zapcore.NewNopCore()),
		usersDB:      usersDB,
		sessionsDB:   sessionsDB,
		tokenCreator: tokenCreator,
		keyGenerator: keyGenerator,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(usersDB, sessionsDB, keyGenerator, tokenCreator)
			got, err := s.RegisterGuestAccount(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterGuestAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegisterGuestAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}
