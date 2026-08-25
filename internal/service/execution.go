package service

import (
	"context"

	"taskqueue/internal/runner"
)

type ExecutionService struct{ dispatcher runner.Dispatcher }

func NewExecutionService(attempts int) *ExecutionService {
	return &ExecutionService{dispatcher: runner.Dispatcher{Attempts: attempts}}
}

func (s *ExecutionService) Run(ctx context.Context, work func(context.Context) error) error {
	return <-s.dispatcher.Execute(context.Background(), work)
}
