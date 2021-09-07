package service

import (
	"auth/mocks"
	"context"
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
		name        string
		args        args
		wantOk      bool
		wantErr     bool
		sessionMock func(sessions *mocks.MockISessions)
	}{
		{
			name: "Success: everything is allright",
			args: args{
				ctx:   context.Background(),
				token: "tokenCreator",
			},
			wantOk:  true,
			wantErr: false,
			sessionMock: func(sessions *mocks.MockISessions) {
				sessions.EXPECT().GetUserIDByToken("tokenCreator")
			},
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
			tt.sessionMock(sessionsDB)
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
