package service

import (
	"auth/internal/infrastructure/db"
	"auth/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func Test_service_CheckToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		token   func(sessions *mocks.MockISessions)
		wantOk  bool
		wantErr bool
	}{
		{
			name: "success: everything is allright",
			args: args{
				ctx:   context.Background(),
				token: "token",
			},
			wantOk:  true,
			wantErr: false,
			token: func(sessions *mocks.MockISessions) {
				sessions.EXPECT().GetUserIDByToken("token").Return(int64(1), nil)
			},
		},
		{
			name: "failed: request to bd error",
			args: args{
				ctx:   context.Background(),
				token: "token",
			},
			token: func(sessions *mocks.MockISessions) {
				sessions.EXPECT().GetUserIDByToken("token").Return(int64(0), errors.New("some error"))
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "failed: request to bd error not found",
			args: args{
				ctx:   context.Background(),
				token: "token",
			},
			token: func(sessions *mocks.MockISessions) {
				sessions.EXPECT().GetUserIDByToken("token").Return(int64(0), db.NotFound)
			},
			wantOk:  false,
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionsDB := mocks.NewMockISessions(ctrl)

	s := &service{
		logger:     zap.New(zapcore.NewNopCore()),
		sessionsDB: sessionsDB,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.token(sessionsDB)
			gotOk, err := s.CheckToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("CheckToken() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
