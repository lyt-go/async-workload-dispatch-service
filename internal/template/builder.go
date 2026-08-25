package template

import (
	"fmt"
	"sync"
)

type Definition struct {
	Name  string
	Steps []string
	Ready bool
}

type Builder struct {
	mu    sync.RWMutex
	cache map[string]*Definition
}

func NewBuilder() *Builder { return &Builder{cache: make(map[string]*Definition)} }

func (b *Builder) Build(name string, steps []string) (result *Definition, err error) {
	result = &Definition{Name: name}
	defer func() {
		if recovered := recover(); recovered != nil {
			// 构建中途失败：丢弃半成品，既不写入缓存也不返回，
			// 避免后续 Get 读到只含部分步骤的 Definition。
			result = nil
			err = fmt.Errorf("build template: %v", recovered)
		}
	}()
	for _, step := range steps {
		if step == "panic" {
			panic("invalid step")
		}
		result.Steps = append(result.Steps, step)
	}
	result.Ready = true
	// 只有在所有步骤都成功、模板就绪后，才把成品写入缓存。
	b.mu.Lock()
	b.cache[name] = result
	b.mu.Unlock()
	return result, nil
}

func (b *Builder) Get(name string) *Definition {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.cache[name]
}
