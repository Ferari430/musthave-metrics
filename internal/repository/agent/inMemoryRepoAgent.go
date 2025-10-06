package repositoryAgent

import (
	"log"
	"maps"
	"sync"

	models "github.com/Ferari430/musthave-metrics/internal/model"
)

type RepositoryAgent struct {
	mu           sync.RWMutex
	metricsStore map[string]*models.Metrics
}

func NewInMemoryAgentDB() *RepositoryAgent {
	return &RepositoryAgent{
		metricsStore: make(map[string]*models.Metrics),
	}
}

func (r *RepositoryAgent) Add(metrics []*models.Metrics) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, metric := range metrics {
		r.metricsStore[metric.ID] = metric
	}
	r.PrintAll()

}

func (r *RepositoryAgent) PrintAll() {
	for name, metric := range r.metricsStore {
		switch metric.MType {
		case "gauge":
			log.Printf("Name: %v, metric: %v", name, *metric.Value)
		case "counter":
			log.Printf("Name: %v, metric: %v", name, *metric.Delta)
		}
	}
	log.Println("----------------------------------")
}

func (r *RepositoryAgent) AllMetrics() map[string]*models.Metrics {
	r.mu.Lock()
	defer r.mu.Unlock()
	metrics := make(map[string]*models.Metrics, len(r.metricsStore))
	maps.Copy(metrics, r.metricsStore)
	return metrics
}
