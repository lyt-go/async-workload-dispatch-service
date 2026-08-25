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
	// policy.Load 现在保证返回已初始化的规则表（即便 names 为空），
	// 因此空配置也能安全检查、登记第一条规则也不会 panic。
	config := policy.Load(names)
	// 默认校验器与 config 共享同一张规则表，使得 Register 写入后
	// Check 能立即读到；用 nil 接口而非 nil 指针，避免 interface-nil 陷阱。
	return &PolicyService{config: config, validator: &requiredPolicy{rules: config.Rules}}
}

func (s *PolicyService) Check(name string) error {
	if s.validator != nil {
		return s.validator.Validate(name)
	}
	return nil
}

func (s *PolicyService) Register(name string) { s.config.Rules[name] = true }
