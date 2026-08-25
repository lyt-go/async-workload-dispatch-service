package service

import (
	"sort"
	"time"

	"taskqueue/internal/model"
	"taskqueue/pkg/idgen"
)

// CreateSubscription 创建订阅关系，校验消费者与队列存在。
func (s *Service) CreateSubscription(input model.Subscription) (*model.Subscription, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetConsumer(input.ConsumerID); err != nil {
		return nil, err
	}
	if _, err := s.store.GetQueue(input.QueueID); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateSubscription(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建订阅 %s", input.ID)
	return &input, nil
}

// GetSubscription 按 ID 获取订阅。
func (s *Service) GetSubscription(id string) (*model.Subscription, error) {
	return s.store.GetSubscription(id)
}

// ListSubscriptions 分页列出订阅。
func (s *Service) ListSubscriptions(filter model.SubscriptionFilter, page, size int) ([]*model.Subscription, int, error) {
	all := s.store.ListSubscriptions()
	matched := make([]*model.Subscription, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Subscription{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateSubscription 更新订阅。
func (s *Service) UpdateSubscription(id string, input model.Subscription) (*model.Subscription, error) {
	existing, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateSubscription(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteSubscription 删除订阅。
func (s *Service) DeleteSubscription(id string) error {
	if err := s.store.DeleteSubscription(id); err != nil {
		return err
	}
	s.log.Infof("删除订阅 %s", id)
	return nil
}
