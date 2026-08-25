package service

import (
	"sort"
	"time"

	"taskqueue/internal/model"
	"taskqueue/pkg/idgen"
)

// CreateQueue 创建任务队列。
func (s *Service) CreateQueue(input model.Queue) (*model.Queue, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetQueueByTopic(input.Topic); err == nil {
		return nil, model.NewValidationError("topic", "主题已存在")
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateQueue(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建队列 %s (%s)", input.Name, input.ID)
	return &input, nil
}

// GetQueue 按 ID 获取队列。
func (s *Service) GetQueue(id string) (*model.Queue, error) {
	return s.store.GetQueue(id)
}

// ListQueues 分页列出队列。
func (s *Service) ListQueues(filter model.QueueFilter, page, size int) ([]*model.Queue, int, error) {
	all := s.store.ListQueues()
	matched := make([]*model.Queue, 0, len(all))
	for _, q := range all {
		if filter.Match(q) {
			matched = append(matched, q)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Queue{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateQueue 更新队列。
func (s *Service) UpdateQueue(id string, input model.Queue) (*model.Queue, error) {
	existing, err := s.store.GetQueue(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Topic = input.Topic
	existing.MaxRetry = input.MaxRetry
	existing.VisibilityTimeout = input.VisibilityTimeout
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateQueue(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// ChangeQueueStatus 变更队列状态（active/paused）。
func (s *Service) ChangeQueueStatus(id, status string) (*model.Queue, error) {
	existing, err := s.store.GetQueue(id)
	if err != nil {
		return nil, err
	}
	if status != model.QueueActive && status != model.QueuePaused {
		return nil, model.NewValidationError("status", "队列状态不合法")
	}
	existing.Status = status
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateQueue(existing); err != nil {
		return nil, err
	}
	s.log.Infof("队列 %s 状态变更为 %s", id, status)
	return existing, nil
}

// DeleteQueue 删除队列。
func (s *Service) DeleteQueue(id string) error {
	if err := s.store.DeleteQueue(id); err != nil {
		return err
	}
	s.log.Infof("删除队列 %s", id)
	return nil
}
