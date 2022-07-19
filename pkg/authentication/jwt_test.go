package authentication

import (
	"testing"
	"time"

	"github.com/qingwave/weave/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	service := NewJWTService("test")

	testCases := []struct {
		name        string
		user        *model.User
		expectedErr bool
	}{
		{
			name:        "user is nil",
			expectedErr: true,
		},
		{
			name:        "create token success",
			user:        &model.User{ID: 1, Name: "someone"},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := service.CreateToken(tc.user)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Empty(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	testCases := []struct {
		name        string
		user        *model.User
		token       string
		expiresAt   int64
		expectedErr bool
	}{
		{
			name:        "invaild token",
			token:       "some-token",
			expectedErr: true,
		},
		{
			name:        "token expiration",
			user:        &model.User{ID: 1, Name: "someone"},
			expiresAt:   int64(-24 * time.Hour),
			expectedErr: true,
		},
		{
			name:        "parse token success",
			user:        &model.User{ID: 1, Name: "someone"},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewJWTService("test")
			if tc.expiresAt != 0 {
				service.expireDuration = tc.expiresAt
			}

			if tc.token == "" {
				token, err := service.CreateToken(tc.user)
				assert.Empty(t, err)
				tc.token = token
			}

			user, err := service.ParseToken(tc.token)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Empty(t, err)
				assert.Equal(t, tc.user.ID, user.ID)
			}
		})
	}
}
