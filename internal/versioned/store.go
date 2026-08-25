package versioned

import "sync"

type Job struct {
	ID      string
	Status  string
	Version int
}

type Store struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewStore() *Store { return &Store{jobs: make(map[string]Job)} }

// Put 写入一条作业，遵循乐观版本号：仅当 job 为新作业，或其版本号大于等于
// 已存在版本时才生效。这样迟到的旧回调（低版本）无法覆盖更新的终态（高版本），
// 防止重试成功后的终态被回退为 running。返回是否真正写入。
func (s *Store) Put(job Job) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.jobs[job.ID]
	if ok && job.Version < cur.Version {
		// 旧版本写入，丢弃，不覆盖更新的状态。
		return false
	}
	s.jobs[job.ID] = job
	return true
}

func (s *Store) Get(id string) Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[id]
}
