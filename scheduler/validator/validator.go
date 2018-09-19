package validator

import (
	"errors"
	"time"

	"github.com/dwarvesf/yggdrasil/toolkit"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
)

// validate request errors
var (
	ErrInvalidService = errors.New("INVALID_SERVICE")
	ErrRequestExpired = errors.New("REQUEST_EXPRIED")
	ErrInvalidRetry   = errors.New("INVALID_RETRY")
)

// ValidateRequest use to validate if a request is valid
func ValidateRequest(r model.Request) error {
	if !isValidService(r.Service) {
		return ErrInvalidService
	}

	if isRequestTimeExpired(r.Timestamp) {
		return ErrRequestExpired
	}

	if !isRetryValid(r.Retry) {
		return ErrInvalidRetry
	}

	return nil
}

func availableSerivces() []string {
	return []string{"email", "sms", "notification", "payment"}
}

func isValidService(service string) bool {
	for _, s := range availableSerivces() {
		if s == service {
			return true
		}
	}

	return false
}

func isRequestTimeExpired(t time.Time) bool {
	return t.Unix() < time.Now().Unix()
}

func isRetryValid(r toolkit.RetryMetadata) bool {
	return r.RetryAfter > 0 && r.CurrenyRetry > 0 && r.MaxRetry >= r.CurrenyRetry
}
