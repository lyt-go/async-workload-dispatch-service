package store

import (
	"taskqueue/internal/model"
)

func (s *MemoryStore) CreateDeadLetter(d *model.DeadLetter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deadLetters[d.ID]; ok {
		return ErrConflict
	}
	s.deadLetters[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDeadLetter(id string) (*model.DeadLetter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.deadLetters[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) GetDeadLetterByTaskID(taskID string) (*model.DeadLetter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, d := range s.deadLetters {
		if d.TaskID == taskID {
			return d, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListDeadLetters() []*model.DeadLetter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DeadLetter, 0, len(s.deadLetters))
	for _, d := range s.deadLetters {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) DeleteDeadLetter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deadLetters[id]; !ok {
		return ErrNotFound
	}
	delete(s.deadLetters, id)
	return nil
}
