package repository

import (
	"fmt"
	"log"
	"sync"

	models "github.com/Ferari430/musthave-metrics/internal/model"
)

type InMemoryRepo struct {
	memStorage map[string]*models.Metrics
	mu         sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	storage := make(map[string]*models.Metrics)

	return &InMemoryRepo{memStorage: storage, mu: sync.RWMutex{}}
}

func (r *InMemoryRepo) Add(metrics *models.Metrics) {
	r.mu.Lock()
	defer r.mu.Unlock()
	log.Printf("Adding/Updating metric: ID=%s, Type=%s", metrics.ID, metrics.MType)
	r.memStorage[metrics.ID] = metrics
	r.PrintAll()
}

func (r *InMemoryRepo) Get(name string) (*models.Metrics, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.memStorage[name]

	return val, ok
}

func (r *InMemoryRepo) Metrics() map[string]*models.Metrics {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.memStorage
}

func (r *InMemoryRepo) PrintAll() {
	log.Println("--- Current state of metrics ---")
	if len(r.memStorage) == 0 {
		log.Println("Storage is empty.")
		return
	}

	for key, metric := range r.memStorage {
		var valueStr string
		if metric.MType == models.Counter && metric.Delta != nil {
			valueStr = fmt.Sprintf("value: %d", *metric.Delta)
		} else if metric.MType == models.Gauge && metric.Value != nil {
			valueStr = fmt.Sprintf("value: %g", *metric.Value)
		} else {
			valueStr = "value: <not set>"
		}
		log.Printf("  - %s (%s) -> %s", key, metric.MType, valueStr)
	}
	log.Println("--------------------------------")

}
