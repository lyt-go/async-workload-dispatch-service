package model

import (
	"strings"
	"time"
)

// Consumer 状态常量。
const (
	ConsumerOnline  = "online"
	ConsumerOffline = "offline"
)

// Consumer 表示一个任务消费者实例。
type Consumer struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	Concurrency   int       `json:"concurrency"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Validate 校验消费者字段。
func (c *Consumer) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return NewValidationError("name", "消费者名称不能为空")
	}
	if c.Concurrency <= 0 {
		return NewValidationError("concurrency", "并发度必须大于 0")
	}
	if c.Status == "" {
		c.Status = ConsumerOnline
	}
	if c.Status != ConsumerOnline && c.Status != ConsumerOffline {
		return NewValidationError("status", "消费者状态不合法")
	}
	return nil
}

// ConsumerFilter 消费者列表筛选条件。
type ConsumerFilter struct {
	Status  string
	Keyword string
}

// Match 判断消费者是否匹配筛选条件。
func (f ConsumerFilter) Match(c *Consumer) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
