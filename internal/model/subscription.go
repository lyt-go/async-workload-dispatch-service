package model

import (
	"strings"
	"time"
)

// Subscription 状态常量。
const (
	SubscriptionActive = "active"
	SubscriptionPaused = "paused"
)

// Subscription 表示消费者与队列之间的订阅关系。
type Subscription struct {
	ID         string    `json:"id"`
	ConsumerID string    `json:"consumer_id"`
	QueueID    string    `json:"queue_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 校验订阅字段。
func (s *Subscription) Validate() error {
	s.ConsumerID = strings.TrimSpace(s.ConsumerID)
	s.QueueID = strings.TrimSpace(s.QueueID)
	if s.ConsumerID == "" {
		return NewValidationError("consumer_id", "消费者 ID 不能为空")
	}
	if s.QueueID == "" {
		return NewValidationError("queue_id", "队列 ID 不能为空")
	}
	if s.Status == "" {
		s.Status = SubscriptionActive
	}
	if s.Status != SubscriptionActive && s.Status != SubscriptionPaused {
		return NewValidationError("status", "订阅状态不合法")
	}
	return nil
}

// SubscriptionFilter 订阅列表筛选条件。
type SubscriptionFilter struct {
	ConsumerID string
	QueueID    string
	Status     string
}

// Match 判断订阅是否匹配筛选条件。
func (f SubscriptionFilter) Match(s *Subscription) bool {
	if f.ConsumerID != "" && s.ConsumerID != f.ConsumerID {
		return false
	}
	if f.QueueID != "" && s.QueueID != f.QueueID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}
