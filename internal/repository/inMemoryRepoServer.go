package repository

import (
	"errors"
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

func (r *InMemoryRepo) MetricJSON(metric *models.Metrics) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	oldmetric, ok := r.memStorage[metric.ID]
	if !ok {
		return false
	}

	// Дополняем только nil-поля
	switch metric.MType {
	case "gauge":
		if metric.Value == nil && oldmetric.Value != nil {
			// создаём новый указатель и копируем значение
			v := *oldmetric.Value
			metric.Value = &v
		}
	case "counter":
		if metric.Delta == nil && oldmetric.Delta != nil {
			d := *oldmetric.Delta
			metric.Delta = &d
		}
	}

	log.Printf("new value: %v, new delta: %v\n", metric.Value, metric.Delta)
	return true
}

func (r *InMemoryRepo) Metrics() map[string]*models.Metrics {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.memStorage
}

func (r *InMemoryRepo) Update(metric *models.Metrics) (*models.Metrics, error) {

	val, ok := r.memStorage[metric.ID]
	log.Println("finded metric in db", ok, val)
	if !ok {
		//logger

		return nil, errors.New("cant find metric in database")
	}

	r.memStorage[val.ID] = metric
	log.Printf("old metric value: %v. Updated metric value %v\n", *val.Value, *metric.Value)
	return metric, nil
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

func (r *InMemoryRepo) Metric(name string) (*models.Metrics, bool) {
	return nil, false
}
