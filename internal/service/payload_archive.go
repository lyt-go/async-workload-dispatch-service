package service

import (
	"sync"

	"taskqueue/internal/payload"
)

type PayloadExporter interface{ Export(id string, value []byte) }

type PayloadArchive struct {
	parser   payload.Parser
	exporter PayloadExporter
	mu       sync.RWMutex
	cache    map[string][]byte
	wg       sync.WaitGroup
}

func NewPayloadArchive(exporter PayloadExporter) *PayloadArchive {
	return &PayloadArchive{exporter: exporter, cache: make(map[string][]byte)}
}

func (a *PayloadArchive) Submit(id, value string) {
	parsed := a.parser.Parse(value)
	a.mu.Lock()
	a.cache[id] = parsed
	a.mu.Unlock()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.exporter.Export(id, parsed)
	}()
}

func (a *PayloadArchive) Get(id string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return string(a.cache[id])
}

func (a *PayloadArchive) Wait() { a.wg.Wait() }
