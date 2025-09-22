package contek_test

import (
	"context"
	"testing"
	"time"

	"github.com/SoeltanIT/agg-common-be/contek"
	"github.com/SoeltanIT/agg-common-be/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetUserContext(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() context.Context
		want    *types.JWTClaims
		wantErr bool
	}{
		{
			name: "success - valid user token",
			setup: func() context.Context {
				claims := &types.JWTClaims{
					ID:        "123",
					Email:     "test@example.com",
					Namespace: "test-ns",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				return context.WithValue(context.Background(), "user", token)
			},
			want: &types.JWTClaims{
				ID:        "123",
				Email:     "test@example.com",
				Namespace: "test-ns",
			},
			wantErr: false,
		},
		{
			name: "error - no user in context",
			setup: func() context.Context {
				return context.Background()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - invalid user type in context",
			setup: func() context.Context {
				return context.WithValue(context.Background(), "user", "not-a-token")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			got := contek.GetUserContext(ctx)

			if tt.wantErr {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.Namespace, got.Namespace)
		})
	}
}

func TestGetUserRawToken(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() context.Context
		want    string
		wantErr bool
	}{
		{
			name: "success - get raw token",
			setup: func() context.Context {
				claims := &types.JWTClaims{
					ID: "123",
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				token.Raw = "raw-token-string"
				return context.WithValue(context.Background(), "user", token)
			},
			want:    "raw-token-string",
			wantErr: false,
		},
		{
			name: "error - no user in context",
			setup: func() context.Context {
				return context.Background()
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("GetUserRawToken() panicked: %v", r)
				}
			}()

			got := contek.GetUserRawToken(ctx)

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
