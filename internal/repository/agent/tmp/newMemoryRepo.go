package repositoryAgent

// import (
// 	"log"
// 	"sync"

// 	models "github.com/Ferari430/musthave-metrics/internal/model"
// )

// type RepositoryAgent struct {
// 	mu      sync.RWMutex
// 	metrics map[string]*models.Metrics
// }

// func NewInMemoryAgentDB() *RepositoryAgent {
// 	return &RepositoryAgent{
// 		metrics: make(map[string]float64),
// 	}
// }

// func (r *RepositoryAgent) Add(metrics map[string]float64) {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	for name, val := range metrics {
// 		newMetric := &models.Metrics{
// 			ID:    name,
// 			MType: "",
// 		}

// 	}

// 	log.Println("Metrics added to repository")

// }

// func (r *RepositoryAgent) GetAllMetrics() map[string]float64 {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
// 	copy := make(map[string]float64)
// 	for k, v := range r.metrics {
// 		copy[k] = v
// 	}
// 	return copy
// }
