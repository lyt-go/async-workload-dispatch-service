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
	b.mu.Lock()
	b.cache[name] = result
	b.mu.Unlock()
	defer func() {
		if recovered := recover(); recovered != nil {
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
	return result, nil
}

func (b *Builder) Get(name string) *Definition {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.cache[name]
}
