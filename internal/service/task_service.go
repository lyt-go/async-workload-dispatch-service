package service

import (
	"sort"
	"time"

	"taskqueue/internal/model"
	"taskqueue/pkg/idgen"
)

// Enqueue 向队列投递一条任务。
func (s *Service) Enqueue(queueID, payload string, priority int) (*model.Task, error) {
	q, err := s.store.GetQueue(queueID)
	if err != nil {
		return nil, err
	}
	if q.Status != model.QueueActive {
		return nil, model.NewValidationError("queue_id", "队列已暂停，无法投递")
	}
	maxAttempts := q.MaxRetry + 1
	t := &model.Task{
		ID:          idgen.HexN(8),
		QueueID:     queueID,
		Payload:     payload,
		Status:      model.TaskPending,
		Priority:    priority,
		Attempts:    0,
		MaxAttempts: maxAttempts,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateTask(t); err != nil {
		return nil, err
	}
	s.log.Infof("投递任务 %s 到队列 %s", t.ID, queueID)
	return t, nil
}

// GetTask 按 ID 获取任务。
func (s *Service) GetTask(id string) (*model.Task, error) {
	return s.store.GetTask(id)
}

// ListTasks 分页列出任务。
func (s *Service) ListTasks(filter model.TaskFilter, page, size int) ([]*model.Task, int, error) {
	all := s.store.ListTasks()
	matched := make([]*model.Task, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Priority != matched[j].Priority {
			return matched[i].Priority > matched[j].Priority
		}
		return matched[i].CreatedAt.Before(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Task{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// Dequeue 从队列取出最高优先级的待执行任务并置为运行中。
func (s *Service) Dequeue(queueID string) (*model.Task, error) {
	q, err := s.store.GetQueue(queueID)
	if err != nil {
		return nil, err
	}
	if q.Status != model.QueueActive {
		return nil, model.NewValidationError("queue_id", "队列已暂停，无法取任务")
	}
	now := time.Now()
	candidates := make([]*model.Task, 0)
	for _, t := range s.store.ListTasks() {
		if t.QueueID == queueID && t.Status == model.TaskPending &&
			(t.NextRunAt.IsZero() || !t.NextRunAt.After(now)) {
			candidates = append(candidates, t)
		}
	}
	if len(candidates) == 0 {
		return nil, model.NewValidationError("queue_id", "队列暂无可用任务")
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority > candidates[j].Priority
		}
		return candidates[i].CreatedAt.Before(candidates[j].CreatedAt)
	})
	t := candidates[0]
	if !model.CanTransitionTask(t.Status, model.TaskRunning) {
		return nil, model.NewValidationError("status", "非法状态流转: "+t.Status+" -> running")
	}
	t.Status = model.TaskRunning
	t.UpdatedAt = now
	if err := s.store.UpdateTask(t); err != nil {
		return nil, err
	}
	s.log.Infof("任务 %s 被取出执行", t.ID)
	return t, nil
}

// AckTask 确认任务执行成功。
func (s *Service) AckTask(id string) (*model.Task, error) {
	t, err := s.store.GetTask(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionTask(t.Status, model.TaskSucceeded) {
		return nil, model.NewValidationError("status", "非法状态流转: "+t.Status+" -> succeeded")
	}
	t.Status = model.TaskSucceeded
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateTask(t); err != nil {
		return nil, err
	}
	s.log.Infof("任务 %s 执行成功", id)
	return t, nil
}

// FailTask 标记任务失败，按剩余次数决定重试或进入死信。
func (s *Service) FailTask(id, reason string) (*model.Task, error) {
	t, err := s.store.GetTask(id)
	if err != nil {
		return nil, err
	}
	if t.Status != model.TaskRunning {
		return nil, model.NewValidationError("status", "只有运行中的任务才能标记失败")
	}
	t.Attempts++
	t.UpdatedAt = time.Now()
	if t.Attempts >= t.MaxAttempts {
		t.Status = model.TaskDead
		if err := s.store.UpdateTask(t); err != nil {
			return nil, err
		}
		dl := &model.DeadLetter{
			ID:              idgen.HexN(8),
			TaskID:          t.ID,
			QueueID:         t.QueueID,
			Reason:          reason,
			OriginalPayload: t.Payload,
			FailedAt:        time.Now(),
			CreatedAt:       time.Now(),
		}
		_ = s.store.CreateDeadLetter(dl)
		s.log.Warnf("任务 %s 进入死信队列，原因：%s", id, reason)
	} else {
		t.Status = model.TaskPending
		t.NextRunAt = time.Now().Add(s.backoff(t.Attempts))
		if err := s.store.UpdateTask(t); err != nil {
			return nil, err
		}
		s.recordRetry(t, reason)
		s.log.Warnf("任务 %s 第 %d 次失败，安排重试", id, t.Attempts)
	}
	return t, nil
}

// RequeueDeadLetter 将死信任务重新投递回待执行状态。
func (s *Service) RequeueDeadLetter(taskID string) (*model.Task, error) {
	t, err := s.store.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionTask(t.Status, model.TaskPending) {
		return nil, model.NewValidationError("status", "非法状态流转: "+t.Status+" -> pending")
	}
	t.Status = model.TaskPending
	t.Attempts = 0
	t.NextRunAt = time.Time{}
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateTask(t); err != nil {
		return nil, err
	}
	if dl, err := s.store.GetDeadLetterByTaskID(taskID); err == nil {
		_ = s.store.DeleteDeadLetter(dl.ID)
	}
	s.log.Infof("死信任务 %s 重新投递", taskID)
	return t, nil
}

// DeleteTask 删除任务。
func (s *Service) DeleteTask(id string) error {
	if err := s.store.DeleteTask(id); err != nil {
		return err
	}
	s.log.Infof("删除任务 %s", id)
	return nil
}

func (s *Service) recordRetry(t *model.Task, reason string) {
	r := &model.RetryRecord{
		ID:          idgen.HexN(8),
		TaskID:      t.ID,
		Attempt:     t.Attempts,
		Error:       reason,
		ScheduledAt: t.NextRunAt,
		CreatedAt:   time.Now(),
	}
	_ = s.store.CreateRetryRecord(r)
}

// backoff 计算指数退避时长（毫秒），上限 60 秒。
func (s *Service) backoff(attempt int) time.Duration {
	base := 1000
	if s.cfg != nil && s.cfg.BackoffBaseMs > 0 {
		base = s.cfg.BackoffBaseMs
	}
	shift := attempt - 1
	if shift > 6 {
		shift = 6
	}
	d := time.Duration(base) * time.Duration(1<<shift) * time.Millisecond
	if d > 60*time.Second {
		d = 60 * time.Second
	}
	return d
}
