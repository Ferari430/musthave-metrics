package repositoryAgent

import (
	"log"
	"sync"
)

type RepositoryAgent struct {
	mu      sync.RWMutex
	metrics map[string]float64
}

func NewInMemoryAgentDB() *RepositoryAgent {
	return &RepositoryAgent{
		metrics: make(map[string]float64),
	}
}

func (r *RepositoryAgent) Add(metrics map[string]float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, v := range metrics {
		r.metrics[k] = v
	}

	log.Println("Metrics added to repository")

}

func (r *RepositoryAgent) GetAllMetrics() map[string]float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(map[string]float64)
	for k, v := range r.metrics {
		copy[k] = v
	}
	return copy
}
