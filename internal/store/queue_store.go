package store

import (
	"taskqueue/internal/model"
)

func (s *MemoryStore) CreateQueue(q *model.Queue) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.queues {
		if exist.Topic == q.Topic {
			return ErrConflict
		}
	}
	s.queues[q.ID] = q
	return nil
}

func (s *MemoryStore) GetQueue(id string) (*model.Queue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.queues[id]
	if !ok {
		return nil, ErrNotFound
	}
	return q, nil
}

func (s *MemoryStore) GetQueueByTopic(topic string) (*model.Queue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, q := range s.queues {
		if q.Topic == topic {
			return q, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListQueues() []*model.Queue {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Queue, 0, len(s.queues))
	for _, q := range s.queues {
		list = append(list, q)
	}
	return list
}

func (s *MemoryStore) UpdateQueue(q *model.Queue) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.queues[q.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.queues {
		if exist.ID != q.ID && exist.Topic == q.Topic {
			return ErrConflict
		}
	}
	s.queues[q.ID] = q
	return nil
}

func (s *MemoryStore) DeleteQueue(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.queues[id]; !ok {
		return ErrNotFound
	}
	delete(s.queues, id)
	return nil
}
