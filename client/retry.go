package client

import (
	"math"
	"time"
)

func shouldRetry(statusCode int) bool {
	return statusCode >= 500 || statusCode == 429
}

func getNextBackoff(attempt int, config RetryConfig) time.Duration {
	multiplier := math.Pow(config.BackoffFactor, float64(attempt))
	backoffDuration := config.InitialInterval * time.Duration(multiplier)
	
	if backoffDuration > config.MaxInterval {
		backoffDuration = config.MaxInterval
	}
	
	return backoffDuration
}
