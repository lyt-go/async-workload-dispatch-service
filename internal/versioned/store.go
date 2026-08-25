package versioned

import "sync"

type Job struct {
	ID      string
	Status  string
	Version int
}

type Store struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewStore() *Store { return &Store{jobs: make(map[string]Job)} }

func (s *Store) Put(job Job) {
	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()
}

func (s *Store) Get(id string) Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[id]
}
