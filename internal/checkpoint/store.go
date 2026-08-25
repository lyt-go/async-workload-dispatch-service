package checkpoint

import "sync"

type Store struct {
	mu     sync.RWMutex
	status map[string]string
}

func NewStore() *Store { return &Store{status: make(map[string]string)} }

type Tx struct {
	store  *Store
	jobID  string
	status string
}

func (s *Store) begin(jobID string) *Tx { return &Tx{store: s, jobID: jobID} }

func (tx *Tx) SetStatus(status string) { tx.status = status }

func (tx *Tx) commit() error {
	tx.store.mu.Lock()
	defer tx.store.mu.Unlock()
	tx.store.status[tx.jobID] = tx.status
	return nil
}

func (tx *Tx) rollback() { tx.status = "" }

func (s *Store) WithTransaction(jobID string, fn func(*Tx) error) (err error) {
	tx := s.begin(jobID)
	defer func() { err = tx.commit() }()
	if err = fn(tx); err != nil {
		tx.rollback()
		return err
	}
	return nil
}

func (s *Store) Status(jobID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status[jobID]
}
