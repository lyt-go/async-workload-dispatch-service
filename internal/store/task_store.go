package store

import (
	"taskqueue/internal/model"
)

func (s *MemoryStore) CreateTask(t *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[t.ID]; ok {
		return ErrConflict
	}
	s.tasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTask(id string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTasks() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTask(t *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[t.ID]; !ok {
		return ErrNotFound
	}
	s.tasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
