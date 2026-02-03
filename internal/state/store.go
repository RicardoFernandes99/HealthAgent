package state

import (
	"sync"

	"github.com/RicardoFernandes99/health-agent/internal/collectors"
)

type store struct {
	mu      sync.RWMutex
	metrics map[string]collectors.Metric
}

func NewStore() *store {
	return &store{
		metrics: make(map[string]collectors.Metric),
	}
}

func (s *store) Set(metric collectors.Metric) {
	s.mu.Lock()         // bloqueia para escrever
	defer s.mu.Unlock() // garante que desbloqueia no fim da execução
	s.metrics[metric.Name] = metric
}

func (s *store) Get(name string) (collectors.Metric, bool) {
	s.mu.RLock() // goroutines podem ler em simultaneo, se estiver um lock ativo para escrever aguarda que termine e só depois podem ler
	defer s.mu.RUnlock()
	metrics, exists := s.metrics[name]

	return metrics, exists
}
