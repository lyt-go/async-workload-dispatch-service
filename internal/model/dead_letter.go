package model

import (
	"strings"
	"time"
)

// DeadLetter 表示一条进入死信队列的任务。
type DeadLetter struct {
	ID              string    `json:"id"`
	TaskID          string    `json:"task_id"`
	QueueID         string    `json:"queue_id"`
	Reason          string    `json:"reason"`
	OriginalPayload string    `json:"original_payload"`
	FailedAt        time.Time `json:"failed_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// Validate 校验死信字段。
func (d *DeadLetter) Validate() error {
	d.TaskID = strings.TrimSpace(d.TaskID)
	d.QueueID = strings.TrimSpace(d.QueueID)
	d.Reason = strings.TrimSpace(d.Reason)
	d.OriginalPayload = strings.TrimSpace(d.OriginalPayload)
	if d.TaskID == "" {
		return NewValidationError("task_id", "任务 ID 不能为空")
	}
	if d.QueueID == "" {
		return NewValidationError("queue_id", "队列 ID 不能为空")
	}
	if d.Reason == "" {
		return NewValidationError("reason", "死信原因不能为空")
	}
	return nil
}

// DeadLetterFilter 死信列表筛选条件。
type DeadLetterFilter struct {
	QueueID string
	Keyword string
}

// Match 判断死信是否匹配筛选条件。
func (f DeadLetterFilter) Match(d *DeadLetter) bool {
	if f.QueueID != "" && d.QueueID != f.QueueID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Reason), k) &&
			!strings.Contains(strings.ToLower(d.OriginalPayload), k) {
			return false
		}
	}
	return true
}
