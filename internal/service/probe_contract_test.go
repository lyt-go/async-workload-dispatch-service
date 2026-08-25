package service_test

import (
	"context"
	"errors"
	"testing"

	"taskqueue/internal/service"
)

func TestProbeUsesFreshContextAfterCanceledRequest(t *testing.T) {
	svc := service.NewProbeService()
	first, cancel := context.WithCancel(context.Background())
	cancel()
	firstCalls := 0
	err := svc.Probe(first, func(ctx context.Context) error {
		firstCalls++
		return ctx.Err()
	})
	if !errors.Is(err, context.Canceled) || firstCalls != 0 {
		t.Fatalf("canceled probe must stop before downstream retry, calls=%d err=%v", firstCalls, err)
	}
	secondCalls := 0
	err = svc.Probe(context.Background(), func(ctx context.Context) error {
		secondCalls++
		return ctx.Err()
	})
	if err != nil || secondCalls != 1 {
		t.Fatalf("fresh probe after cancellation must use its own context once, calls=%d err=%v", secondCalls, err)
	}
}
