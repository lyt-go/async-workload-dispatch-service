package store

import (
	"taskqueue/internal/model"
)

func (s *MemoryStore) CreateRetryRecord(r *model.RetryRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retries[r.ID]; ok {
		return ErrConflict
	}
	s.retries[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRetryRecord(id string) (*model.RetryRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.retries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRetryRecords() []*model.RetryRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RetryRecord, 0, len(s.retries))
	for _, r := range s.retries {
		list = append(list, r)
	}
	return list
}
