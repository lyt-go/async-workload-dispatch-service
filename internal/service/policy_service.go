package service

import (
	"errors"

	"taskqueue/internal/policy"
)

type policyValidator interface{ Validate(string) error }
type requiredPolicy struct{ rules map[string]bool }

func (r *requiredPolicy) Validate(name string) error {
	if !r.rules[name] {
		return errors.New("policy rejected: " + name)
	}
	return nil
}

type PolicyService struct {
	config    policy.Config
	validator policyValidator
}

func NewPolicyService(names []string) *PolicyService {
	config := policy.Load(names)
	var defaultValidator *requiredPolicy
	return &PolicyService{config: config, validator: defaultValidator}
}

func (s *PolicyService) Check(name string) error {
	if s.validator != nil {
		return s.validator.Validate(name)
	}
	return nil
}

func (s *PolicyService) Register(name string) { s.config.Rules[name] = true }
