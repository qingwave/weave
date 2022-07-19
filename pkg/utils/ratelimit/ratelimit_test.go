package ratelimit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      *LimitConfig
		expectedErr error
	}{
		{
			name:        "invalid burst and qps",
			config:      &LimitConfig{},
			expectedErr: errors.New("LimitConfig Burst and QPS cannot be empty"),
		},
		{
			name:        "qps great than burst",
			config:      &LimitConfig{QPS: 10, Burst: 1},
			expectedErr: errors.New("LimitConfig QPS(10) must less than Burst(1)"),
		},
		{
			name:        "set cache size",
			config:      &LimitConfig{QPS: 2, Burst: 4, CacheSize: 100},
			expectedErr: nil,
		},
		{
			name:        "config success",
			config:      &LimitConfig{QPS: 2, Burst: 4},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cacheSize := tc.config.CacheSize
			err := tc.config.Validate()
			assert.Equal(t, tc.expectedErr, err)

			if cacheSize == 0 && tc.expectedErr == nil {
				assert.Equal(t, defaultCacheSize, tc.config.CacheSize)
			}
		})
	}
}

func TestRateLimiter(t *testing.T) {
	testCases := []struct {
		name               string
		config             *LimitConfig
		currentRquestCount int
		shouldAccepted     bool
		expectedErr        bool
	}{
		{
			name:        "null config",
			config:      nil,
			expectedErr: true,
		},
		{
			name:        "invaild type",
			config:      &LimitConfig{LimitType: "some-type"},
			expectedErr: true,
		},
		{
			name:           "accept request",
			config:         &LimitConfig{QPS: 1, Burst: 1, LimitType: ServerLimitType},
			shouldAccepted: true,
		},
		{
			name:               "reject request",
			config:             &LimitConfig{QPS: 1, Burst: 1, LimitType: ServerLimitType},
			currentRquestCount: 10,
			shouldAccepted:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewRateLimiter(tc.config)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}

			if tc.currentRquestCount > 0 {
				for i := 0; i < tc.currentRquestCount; i++ {
					r.Accept(nil)
				}
			}
			err = r.Accept(nil)
			if tc.shouldAccepted {
				assert.Empty(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
