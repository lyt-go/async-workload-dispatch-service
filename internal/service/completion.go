package service

import "taskqueue/internal/checkpoint"

type CompletionPublisher interface {
	Publish(jobID, status string) error
}

type CompletionService struct {
	store     *checkpoint.Store
	publisher CompletionPublisher
}

func NewCompletionService(store *checkpoint.Store, publisher CompletionPublisher) *CompletionService {
	return &CompletionService{store: store, publisher: publisher}
}

func (s *CompletionService) Finalize(jobID string, validate func() error) error {
	return s.store.WithTransaction(jobID, func(tx *checkpoint.Tx) error {
		tx.SetStatus("completed")
		if err := s.publisher.Publish(jobID, "completed"); err != nil {
			return err
		}
		return validate()
	})
}
