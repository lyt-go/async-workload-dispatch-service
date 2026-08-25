package probe

import (
	"context"
	"sync"
)

type ContextStore struct {
	mu  sync.Mutex
	ctx context.Context
}

func (s *ContextStore) ForRequest(ctx context.Context) context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx == nil {
		s.ctx = ctx
	}
	return s.ctx
}
