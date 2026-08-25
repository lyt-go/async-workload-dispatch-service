// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"taskqueue/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Queue
	CreateQueue(q *model.Queue) error
	GetQueue(id string) (*model.Queue, error)
	GetQueueByTopic(topic string) (*model.Queue, error)
	ListQueues() []*model.Queue
	UpdateQueue(q *model.Queue) error
	DeleteQueue(id string) error

	// Task
	CreateTask(t *model.Task) error
	GetTask(id string) (*model.Task, error)
	ListTasks() []*model.Task
	UpdateTask(t *model.Task) error
	DeleteTask(id string) error

	// Consumer
	CreateConsumer(c *model.Consumer) error
	GetConsumer(id string) (*model.Consumer, error)
	ListConsumers() []*model.Consumer
	UpdateConsumer(c *model.Consumer) error
	DeleteConsumer(id string) error

	// Subscription
	CreateSubscription(s *model.Subscription) error
	GetSubscription(id string) (*model.Subscription, error)
	ListSubscriptions() []*model.Subscription
	UpdateSubscription(s *model.Subscription) error
	DeleteSubscription(id string) error

	// RetryRecord
	CreateRetryRecord(r *model.RetryRecord) error
	GetRetryRecord(id string) (*model.RetryRecord, error)
	ListRetryRecords() []*model.RetryRecord

	// DeadLetter
	CreateDeadLetter(d *model.DeadLetter) error
	GetDeadLetter(id string) (*model.DeadLetter, error)
	GetDeadLetterByTaskID(taskID string) (*model.DeadLetter, error)
	ListDeadLetters() []*model.DeadLetter
	DeleteDeadLetter(id string) error
}
