package toolkit

import (
	"time"
)

// RetryMetadata is a struct define retry metadata
type RetryMetadata struct {
	RetryAfter   time.Duration `json:"retryAfter"`
	CurrenyRetry int           `json:"currentRetry"`
	MaxRetry     int           `json:"maxRetry"`
}
