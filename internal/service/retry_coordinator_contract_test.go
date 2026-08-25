package service_test

import (
	"sync"
	"testing"

	"taskqueue/internal/service"
	"taskqueue/internal/versioned"
)

type deduplicatingEffector struct {
	mu   sync.Mutex
	keys map[string]bool
}

func (e *deduplicatingEffector) Apply(key string) {
	e.mu.Lock()
	e.keys[key] = true
	e.mu.Unlock()
}

func TestLateFirstAttemptCannotOverwriteRetrySuccess(t *testing.T) {
	st := versioned.NewStore()
	effect := &deduplicatingEffector{keys: make(map[string]bool)}
	coordinator := service.NewRetryCoordinator(st, effect)
	releaseLate := make(chan struct{})
	close(releaseLate)
	coordinator.CompleteAfterRetry("job-88", releaseLate)
	job := st.Get("job-88")
	if job.Status != "succeeded" || job.Version != 2 {
		t.Fatalf("late first attempt must not overwrite retry success, got status=%s version=%d", job.Status, job.Version)
	}
	effect.mu.Lock()
	count := len(effect.keys)
	effect.mu.Unlock()
	if count != 1 {
		t.Fatalf("retry attempts must share one idempotency effect key, got %d effects", count)
	}
}
