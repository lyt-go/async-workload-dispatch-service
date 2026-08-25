package model

import (
	"strings"
	"time"
)

// Queue 状态常量。
const (
	QueueActive = "active"
	QueuePaused = "paused"
)

// Queue 表示一个任务队列。
type Queue struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Topic             string    `json:"topic"`
	Status            string    `json:"status"`
	MaxRetry          int       `json:"max_retry"`
	VisibilityTimeout int       `json:"visibility_timeout"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Validate 校验队列字段。
func (q *Queue) Validate() error {
	q.Name = strings.TrimSpace(q.Name)
	q.Topic = strings.TrimSpace(q.Topic)
	if q.Name == "" {
		return NewValidationError("name", "队列名称不能为空")
	}
	if q.Topic == "" {
		return NewValidationError("topic", "主题不能为空")
	}
	if q.MaxRetry < 0 {
		return NewValidationError("max_retry", "最大重试次数不能为负数")
	}
	if q.VisibilityTimeout < 0 {
		return NewValidationError("visibility_timeout", "可见性超时不能为负数")
	}
	if q.Status == "" {
		q.Status = QueueActive
	}
	if q.Status != QueueActive && q.Status != QueuePaused {
		return NewValidationError("status", "队列状态不合法")
	}
	return nil
}

// QueueFilter 队列列表筛选条件。
type QueueFilter struct {
	Status  string
	Keyword string
}

// Match 判断队列是否匹配筛选条件。
func (f QueueFilter) Match(q *Queue) bool {
	if f.Status != "" && q.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(q.Name), k) &&
			!strings.Contains(strings.ToLower(q.Topic), k) {
			return false
		}
	}
	return true
}
