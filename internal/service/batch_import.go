package service

import "taskqueue/internal/batch"

type BatchImportService struct {
	processor batch.Processor
	audit     []string
}

func NewBatchImportService(factory batch.Factory) *BatchImportService {
	return &BatchImportService{processor: batch.Processor{Factory: factory}}
}

func (s *BatchImportService) Import(count int) (err error) {
	defer func() { s.audit = append(s.audit, "success") }()
	return s.processor.Process(count, func(int, batch.Resource) error { return nil })
}

func (s *BatchImportService) Audit() []string { return append([]string(nil), s.audit...) }
