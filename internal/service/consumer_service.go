package service

import (
	"sort"
	"time"

	"taskqueue/internal/model"
	"taskqueue/pkg/idgen"
)

// CreateConsumer 创建消费者实例。
func (s *Service) CreateConsumer(input model.Consumer) (*model.Consumer, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	input.LastHeartbeat = now
	if err := s.store.CreateConsumer(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建消费者 %s (%s)", input.Name, input.ID)
	return &input, nil
}

// GetConsumer 按 ID 获取消费者。
func (s *Service) GetConsumer(id string) (*model.Consumer, error) {
	return s.store.GetConsumer(id)
}

// ListConsumers 分页列出消费者。
func (s *Service) ListConsumers(filter model.ConsumerFilter, page, size int) ([]*model.Consumer, int, error) {
	all := s.store.ListConsumers()
	matched := make([]*model.Consumer, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Consumer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateConsumer 更新消费者。
func (s *Service) UpdateConsumer(id string, input model.Consumer) (*model.Consumer, error) {
	existing, err := s.store.GetConsumer(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Concurrency = input.Concurrency
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateConsumer(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Heartbeat 更新消费者心跳时间。
func (s *Service) Heartbeat(id string) (*model.Consumer, error) {
	existing, err := s.store.GetConsumer(id)
	if err != nil {
		return nil, err
	}
	existing.LastHeartbeat = time.Now()
	existing.Status = model.ConsumerOnline
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateConsumer(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteConsumer 删除消费者。
func (s *Service) DeleteConsumer(id string) error {
	if err := s.store.DeleteConsumer(id); err != nil {
		return err
	}
	s.log.Infof("删除消费者 %s", id)
	return nil
}
