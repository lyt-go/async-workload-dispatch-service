package service

import (
	"errors"

	templatepkg "taskqueue/internal/template"
)

type TemplateService struct{ builder *templatepkg.Builder }

func NewTemplateService() *TemplateService {
	return &TemplateService{builder: templatepkg.NewBuilder()}
}

func (s *TemplateService) Create(name string, steps []string) error {
	_, err := s.builder.Build(name, steps)
	return err
}

func (s *TemplateService) StepCount(name string) (int, error) {
	definition := s.builder.Get(name)
	if definition == nil {
		return 0, errors.New("template not found")
	}
	if !definition.Ready {
		return 0, errors.New("template not ready")
	}
	return len(definition.Steps), nil
}
