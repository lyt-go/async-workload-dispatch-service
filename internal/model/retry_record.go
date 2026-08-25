package model

import (
	"strings"
	"time"
)

// RetryRecord 表示一次任务重试记录。
type RetryRecord struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	Attempt     int       `json:"attempt"`
	Error       string    `json:"error"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate 校验重试记录字段。
func (r *RetryRecord) Validate() error {
	r.TaskID = strings.TrimSpace(r.TaskID)
	r.Error = strings.TrimSpace(r.Error)
	if r.TaskID == "" {
		return NewValidationError("task_id", "任务 ID 不能为空")
	}
	if r.Attempt < 1 {
		return NewValidationError("attempt", "重试序号必须大于 0")
	}
	if r.Error == "" {
		return NewValidationError("error", "失败原因不能为空")
	}
	return nil
}

// RetryRecordFilter 重试记录筛选条件。
type RetryRecordFilter struct {
	TaskID string
}

// Match 判断重试记录是否匹配筛选条件。
func (f RetryRecordFilter) Match(r *RetryRecord) bool {
	if f.TaskID != "" && r.TaskID != f.TaskID {
		return false
	}
	return true
}
