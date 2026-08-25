package service_test

import (
	"testing"

	"taskqueue/internal/service"
)

func didPanic(fn func()) (panicked bool) {
	defer func() { panicked = recover() != nil }()
	fn()
	return false
}

func TestEmptyPolicySupportsCheckThenFirstRegistration(t *testing.T) {
	svc := service.NewPolicyService(nil)
	if didPanic(func() { _ = svc.Check("dispatch") }) {
		t.Fatalf("empty policy check must not panic through a typed nil validator")
	}
	if didPanic(func() { svc.Register("dispatch") }) {
		t.Fatalf("first policy registration must not panic on an uninitialized rule map")
	}
	if err := svc.Check("dispatch"); err != nil {
		t.Fatalf("registered dispatch policy must pass validation, got %v", err)
	}
}
