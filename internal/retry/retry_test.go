package retry

import (
	"testing"
	"time"
)

func Test_shouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"should retry on 500", 500, true},
		{"should retry on 502", 502, true},
		{"should retry on 503", 503, true},
		{"should retry on 504", 504, true},
		{"should retry on 429", 429, true},
		{"should not retry on 400", 400, false},
		{"should not retry on 401", 401, false},
		{"should not retry on 403", 403, false},
		{"should not retry on 404", 404, false},
		{"should not retry on 200", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldRetry(tt.statusCode); got != tt.want {
				t.Errorf("shouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNextBackoff(t *testing.T) {
	tests := []struct {
		name           string
		attempt        int
		backoffFactor  float64
		initialInterval time.Duration
		maxInterval    time.Duration
		wantDelay      time.Duration
	}{
		{
			name:            "first attempt with default config",
			attempt:         0,
			backoffFactor:   2.0,
			initialInterval: time.Second,
			maxInterval:     time.Second * 30,
			wantDelay:      time.Second,
		},
		{
			name:            "second attempt with default config",
			attempt:         1,
			backoffFactor:   2.0,
			initialInterval: time.Second,
			maxInterval:     time.Second * 30,
			wantDelay:      time.Second * 2,
		},
		{
			name:            "max interval reached",
			attempt:         10,
			backoffFactor:   2.0,
			initialInterval: time.Second,
			maxInterval:     time.Second * 30,
			wantDelay:      time.Second * 30,
		},
		{
			name:            "custom config with different factor",
			attempt:         2,
			backoffFactor:   1.5,
			initialInterval: time.Second * 2,
			maxInterval:     time.Second * 60,
			wantDelay:      time.Second * 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNextBackoff(tt.attempt, tt.backoffFactor, tt.initialInterval, tt.maxInterval)
			if got != tt.wantDelay {
				t.Errorf("getNextBackoff() = %v, want %v", got, tt.wantDelay)
			}
		})
	}
}
