package snapshot

import "sync"

type Store struct {
	mu    sync.RWMutex
	items map[string]int
}

func NewStore() *Store { return &Store{items: make(map[string]int)} }

func (s *Store) Put(key string, value int) {
	s.mu.Lock()
	s.items[key] = value
	s.mu.Unlock()
}

func (s *Store) Snapshot() map[string]int { return s.items }
