package toolkit

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fatih/structs"
)

// validate retry errors
var (
	ErrInvalidRetry = errors.New("INVALID_RETRY")
)

// SchedulerRequest is a request for scheduler
type SchedulerRequest struct {
	Service   string                 `json:"service"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
	Retry     RetryMetadata          `json:"retry"`
}

// RetryMetadata is a struct define retry metadata
type RetryMetadata struct {
	RetryAfter   time.Duration `json:"retryAfter"`
	CurrenyRetry int           `json:"currentRetry"`
	MaxRetry     int           `json:"maxRetry"`
}

// CreateRetryMessage return retry message
func CreateRetryMessage(service string, payload interface{}, retry RetryMetadata) ([]byte, error) {
	if !retry.IsValidRetry() {
		return nil, ErrInvalidRetry
	}

	req := SchedulerRequest{
		Service:   service,
		Payload:   structs.Map(payload),
		Timestamp: time.Now().Add(retry.RetryAfter),
		Retry:     retry,
	}

	message, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// IsValidRetry return if retry info is valid or not
func (r RetryMetadata) IsValidRetry() bool {
	return r.RetryAfter > 0 && r.CurrenyRetry > 0 && r.MaxRetry >= r.CurrenyRetry
}
