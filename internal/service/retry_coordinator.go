package service

import (
	"sync"

	"taskqueue/internal/versioned"
)

type RetryEffector interface{ Apply(key string) }

type RetryCoordinator struct {
	store  *versioned.Store
	effect RetryEffector
}

func NewRetryCoordinator(store *versioned.Store, effect RetryEffector) *RetryCoordinator {
	return &RetryCoordinator{store: store, effect: effect}
}

// CompleteAfterRetry 模拟一次重试链：首次处理迟返回，重试已成功后旧回调才到达。
// 约束：
//   - 终态不倒退：重试成功写入的 succeeded 不得被迟到的旧回调覆盖回 running。
//   - 同一作业的副作用只生效一次：两次尝试必须使用相同的幂等标识，由 Apply 去重。
func (c *RetryCoordinator) CompleteAfterRetry(jobID string, releaseLate <-chan struct{}) {
	// 用作业自身作为幂等键：无论重试多少次，同一作业的外部副作用只生效一次。
	effectKey := "job:" + jobID

	// 首次尝试：置为 running(v1)，副作用以幂等键提交。
	c.store.Put(versioned.Job{ID: jobID, Status: "running", Version: 1})
	c.effect.Apply(effectKey)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-releaseLate
		// 迟到的旧回调：仍试图写 running(v1)。此时终态已为 succeeded(v2)，
		// 由于版本号更低，Put 会被拒绝，终态不再倒退。
		c.store.Put(versioned.Job{ID: jobID, Status: "running", Version: 1})
	}()

	// 重试成功：写入终态 succeeded(v2)。版本更高，旧回调此后无法覆盖。
	c.store.Put(versioned.Job{ID: jobID, Status: "succeeded", Version: 2})
	// 同一作业的副作用：幂等键不变，重复 Apply 不会二次执行外部动作。
	c.effect.Apply(effectKey)
	wg.Wait()
}
