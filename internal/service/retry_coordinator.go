package service

import (
	"sync"

	"taskqueue/internal/versioned"
)

type RetryEffector interface{ Apply(key string) }

type RetryCoordinator struct {
	store  *versioned.Store
	effect RetryEffector
}

func NewRetryCoordinator(store *versioned.Store, effect RetryEffector) *RetryCoordinator {
	return &RetryCoordinator{store: store, effect: effect}
}

func (c *RetryCoordinator) CompleteAfterRetry(jobID string, releaseLate <-chan struct{}) {
	c.store.Put(versioned.Job{ID: jobID, Status: "running", Version: 1})
	c.effect.Apply(jobID + "-attempt-1")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-releaseLate
		c.store.Put(versioned.Job{ID: jobID, Status: "running", Version: 1})
	}()
	c.store.Put(versioned.Job{ID: jobID, Status: "succeeded", Version: 2})
	c.effect.Apply(jobID + "-attempt-2")
	wg.Wait()
}
