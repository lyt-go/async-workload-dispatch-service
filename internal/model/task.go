package model

import (
	"strings"
	"time"
)

// Task 状态常量。
const (
	TaskPending   = "pending"
	TaskRunning   = "running"
	TaskSucceeded = "succeeded"
	TaskFailed    = "failed"
	TaskDead      = "dead"
)

// taskTransitions 定义任务合法状态流转。
var taskTransitions = map[string]map[string]bool{
	TaskPending: {TaskRunning: true},
	TaskRunning: {TaskSucceeded: true, TaskFailed: true},
	TaskFailed:  {TaskPending: true, TaskDead: true},
	TaskDead:    {TaskPending: true},
}

// CanTransitionTask 判断任务状态能否从 from 流转到 to。
func CanTransitionTask(from, to string) bool {
	if m, ok := taskTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Task 表示队列中的一条任务。
type Task struct {
	ID          string    `json:"id"`
	QueueID     string    `json:"queue_id"`
	Payload     string    `json:"payload"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
	NextRunAt   time.Time `json:"next_run_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验任务字段。
func (t *Task) Validate() error {
	t.QueueID = strings.TrimSpace(t.QueueID)
	t.Payload = strings.TrimSpace(t.Payload)
	if t.QueueID == "" {
		return NewValidationError("queue_id", "队列 ID 不能为空")
	}
	if t.Payload == "" {
		return NewValidationError("payload", "任务负载不能为空")
	}
	if t.Priority < 0 {
		return NewValidationError("priority", "优先级不能为负数")
	}
	if t.MaxAttempts <= 0 {
		return NewValidationError("max_attempts", "最大尝试次数必须大于 0")
	}
	if t.Attempts < 0 {
		return NewValidationError("attempts", "尝试次数不能为负数")
	}
	if t.Status == "" {
		t.Status = TaskPending
	}
	switch t.Status {
	case TaskPending, TaskRunning, TaskSucceeded, TaskFailed, TaskDead:
	default:
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

// TaskFilter 任务列表筛选条件。
type TaskFilter struct {
	QueueID string
	Status  string
}

// Match 判断任务是否匹配筛选条件。
func (f TaskFilter) Match(t *Task) bool {
	if f.QueueID != "" && t.QueueID != f.QueueID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	return true
}
