package service

import "taskqueue/internal/snapshot"

type SnapshotTotalService struct{ store *snapshot.Store }

func NewSnapshotTotalService(store *snapshot.Store) *SnapshotTotalService {
	return &SnapshotTotalService{store: store}
}

func (s *SnapshotTotalService) TotalAsync(ready chan<- struct{}, start <-chan struct{}) <-chan int {
	result := make(chan int, 1)
	view := s.store.Snapshot()
	go func() {
		close(ready)
		<-start
		total := 0
		for _, value := range view {
			total += value
		}
		result <- total
	}()
	return result
}
