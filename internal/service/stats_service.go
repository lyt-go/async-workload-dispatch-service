package service

import (
	"taskqueue/internal/model"
)

// QueueStat 单个队列的任务统计。
type QueueStat struct {
	Name        string `json:"name"`
	Topic       string `json:"topic"`
	Status      string `json:"status"`
	Pending     int    `json:"pending"`
	Running     int    `json:"running"`
	Succeeded   int    `json:"succeeded"`
	Failed      int    `json:"failed"`
	Dead        int    `json:"dead"`
	RetryCount  int    `json:"retry_count"`
	DeadLetters int    `json:"dead_letters"`
}

// QueueStats 队列全局统计。
type QueueStats struct {
	QueueCount      int         `json:"queue_count"`
	ConsumerCount   int         `json:"consumer_count"`
	SubscriptionCnt int         `json:"subscription_count"`
	TaskCount       int         `json:"task_count"`
	PendingCount    int         `json:"pending_count"`
	RunningCount    int         `json:"running_count"`
	SucceededCount  int         `json:"succeeded_count"`
	FailedCount     int         `json:"failed_count"`
	DeadCount       int         `json:"dead_count"`
	RetryCount      int         `json:"retry_count"`
	DeadLetterCount int         `json:"dead_letter_count"`
	ByQueue         []QueueStat `json:"by_queue"`
}

// Overview 汇总任务队列全局统计。
func (s *Service) Overview() (*QueueStats, error) {
	stats := &QueueStats{ByQueue: []QueueStat{}}
	stats.QueueCount = len(s.store.ListQueues())
	stats.ConsumerCount = len(s.store.ListConsumers())
	stats.SubscriptionCnt = len(s.store.ListSubscriptions())
	stats.RetryCount = len(s.store.ListRetryRecords())
	stats.DeadLetterCount = len(s.store.ListDeadLetters())

	byQueue := make(map[string]*QueueStat)
	for _, q := range s.store.ListQueues() {
		byQueue[q.ID] = &QueueStat{Name: q.Name, Topic: q.Topic, Status: q.Status}
	}
	for _, t := range s.store.ListTasks() {
		stats.TaskCount++
		switch t.Status {
		case model.TaskPending:
			stats.PendingCount++
		case model.TaskRunning:
			stats.RunningCount++
		case model.TaskSucceeded:
			stats.SucceededCount++
		case model.TaskFailed:
			stats.FailedCount++
		case model.TaskDead:
			stats.DeadCount++
		}
		if qs, ok := byQueue[t.QueueID]; ok {
			switch t.Status {
			case model.TaskPending:
				qs.Pending++
			case model.TaskRunning:
				qs.Running++
			case model.TaskSucceeded:
				qs.Succeeded++
			case model.TaskFailed:
				qs.Failed++
			case model.TaskDead:
				qs.Dead++
			}
		}
	}
	for _, d := range s.store.ListDeadLetters() {
		if qs, ok := byQueue[d.QueueID]; ok {
			qs.DeadLetters++
		}
	}
	for _, q := range s.store.ListQueues() {
		if qs, ok := byQueue[q.ID]; ok {
			stats.ByQueue = append(stats.ByQueue, *qs)
		}
	}
	return stats, nil
}
