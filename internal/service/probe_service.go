package service

import (
	"context"

	"taskqueue/internal/probe"
)

type ProbeService struct{ contexts *probe.ContextStore }

func NewProbeService() *ProbeService { return &ProbeService{contexts: &probe.ContextStore{}} }

func (s *ProbeService) Probe(ctx context.Context, call func(context.Context) error) error {
	requestCtx := s.contexts.ForRequest(ctx)
	if err := call(requestCtx); err != nil {
		return call(context.Background())
	}
	return nil
}
