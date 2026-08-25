package store

import (
	"taskqueue/internal/model"
)

func (s *MemoryStore) CreateSubscription(sub *model.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.subscriptions {
		if exist.ConsumerID == sub.ConsumerID && exist.QueueID == sub.QueueID {
			return ErrConflict
		}
	}
	s.subscriptions[sub.ID] = sub
	return nil
}

func (s *MemoryStore) GetSubscription(id string) (*model.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.subscriptions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sub, nil
}

func (s *MemoryStore) ListSubscriptions() []*model.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Subscription, 0, len(s.subscriptions))
	for _, sub := range s.subscriptions {
		list = append(list, sub)
	}
	return list
}

func (s *MemoryStore) UpdateSubscription(sub *model.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[sub.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.subscriptions {
		if exist.ID != sub.ID && exist.ConsumerID == sub.ConsumerID && exist.QueueID == sub.QueueID {
			return ErrConflict
		}
	}
	s.subscriptions[sub.ID] = sub
	return nil
}

func (s *MemoryStore) DeleteSubscription(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[id]; !ok {
		return ErrNotFound
	}
	delete(s.subscriptions, id)
	return nil
}
