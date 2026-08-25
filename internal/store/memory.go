package store

import (
	"sync"

	"taskqueue/internal/model"
)

// MemoryStore 基于内存 map 的线程安全存储实现。
type MemoryStore struct {
	mu            sync.RWMutex
	queues        map[string]*model.Queue
	tasks         map[string]*model.Task
	consumers     map[string]*model.Consumer
	subscriptions map[string]*model.Subscription
	retries       map[string]*model.RetryRecord
	deadLetters   map[string]*model.DeadLetter
}

// NewMemoryStore 创建空的 MemoryStore。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		queues:        make(map[string]*model.Queue),
		tasks:         make(map[string]*model.Task),
		consumers:     make(map[string]*model.Consumer),
		subscriptions: make(map[string]*model.Subscription),
		retries:       make(map[string]*model.RetryRecord),
		deadLetters:   make(map[string]*model.DeadLetter),
	}
}

var _ Store = (*MemoryStore)(nil)
