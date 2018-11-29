package validator

import (
	"testing"
	"time"

	"github.com/dwarvesf/yggdrasil/services/scheduler/model"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func TestValidateRequestWhenInvalidService(t *testing.T) {
	r := model.Request{
		Service:   "test",
		Timestamp: time.Now().Add(time.Millisecond * 10),
	}

	if err := ValidateRequest(r); err != ErrInvalidService {
		t.Errorf("Expect error to be ErrInvalidService, but got %v", err)
	}
}

func TestValidateRequestWhenRequestExpired(t *testing.T) {
	r := model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(time.Second * -1),
	}

	if err := ValidateRequest(r); err != ErrRequestExpired {
		t.Errorf("Expect error to be ErrRequestExpired, but got %v", err)
	}
}

func TestValidateRequestWhenRetryWhenInvalid(t *testing.T) {
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 1,
		MaxRetry:     3,
		RetryAfter:   -10 * time.Second,
	}
	r := model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(time.Second * 10),
		Retry:     retry,
	}

	if err := ValidateRequest(r); err != ErrInvalidRetry {
		t.Errorf("Expect error to be ErrInvalidRetry, but got %v", err)
	}
}

func TestValidateRequestWhenRetryWhenCurrentRetryInvalid(t *testing.T) {
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 0,
		MaxRetry:     3,
		RetryAfter:   time.Second,
	}
	r := model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(time.Second * 10),
		Retry:     retry,
	}

	if err := ValidateRequest(r); err != ErrInvalidRetry {
		t.Errorf("Expect error to be ErrInvalidRetry, but got %v", err)
	}
}

func TestValidateRequestWhenRetryWhenMaxRetryInvalid(t *testing.T) {
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 2,
		MaxRetry:     1,
		RetryAfter:   time.Second,
	}
	r := model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(time.Second * 10),
		Retry:     retry,
	}

	if err := ValidateRequest(r); err != ErrInvalidRetry {
		t.Errorf("Expect error to be ErrInvalidRetry, but got %v", err)
	}
}

func TestValidate(t *testing.T) {
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 2,
		MaxRetry:     3,
		RetryAfter:   time.Second,
	}

	r := model.Request{
		Service:   "email",
		Timestamp: time.Now().Add(time.Millisecond * 10),
		Retry:     retry,
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}

	r = model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(time.Millisecond * 10),
		Retry:     retry,
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}

	r = model.Request{
		Service:   "notification",
		Timestamp: time.Now().Add(time.Millisecond * 10),
		Retry:     retry,
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}

	r = model.Request{
		Service:   "payment",
		Timestamp: time.Now().Add(time.Millisecond * 10),
		Retry:     retry,
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}
}
