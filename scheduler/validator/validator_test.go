package validator

import (
	"testing"
	"time"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
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

func TestValidate(t *testing.T) {
	r := model.Request{
		Service:   "email",
		Timestamp: time.Now().Add(time.Millisecond * 10),
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}

	r = model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(time.Millisecond * 10),
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}

	r = model.Request{
		Service:   "notification",
		Timestamp: time.Now().Add(time.Millisecond * 10),
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}

	r = model.Request{
		Service:   "payment",
		Timestamp: time.Now().Add(time.Millisecond * 10),
	}
	if err := ValidateRequest(r); err != nil {
		t.Errorf("Expect error to be nil, but got %v", err)
	}
}
