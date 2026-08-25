package store

import (
	"taskqueue/internal/model"
)

func (s *MemoryStore) CreateConsumer(c *model.Consumer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.consumers[c.ID]; ok {
		return ErrConflict
	}
	s.consumers[c.ID] = c
	return nil
}

func (s *MemoryStore) GetConsumer(id string) (*model.Consumer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.consumers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListConsumers() []*model.Consumer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Consumer, 0, len(s.consumers))
	for _, c := range s.consumers {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateConsumer(c *model.Consumer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.consumers[c.ID]; !ok {
		return ErrNotFound
	}
	s.consumers[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteConsumer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.consumers[id]; !ok {
		return ErrNotFound
	}
	delete(s.consumers, id)
	return nil
}
