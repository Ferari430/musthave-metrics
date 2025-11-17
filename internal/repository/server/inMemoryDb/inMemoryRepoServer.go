package inMemoryStorage

import (
	"errors"
	"fmt"
	"log"
	"sync"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
)

type InMemoryRepo struct {
	memStorage map[string]*models.Metrics
	mu         sync.RWMutex
	logger     *logger.Logger
}

func NewInMemoryStorage(logger *logger.Logger) *InMemoryRepo {
	storage := make(map[string]*models.Metrics)
	return &InMemoryRepo{memStorage: storage, mu: sync.RWMutex{},
		logger: logger,
	}
}

func (r *InMemoryRepo) GetAll() ([]*models.Metrics, error) {
	log.Println("--- Current state of metrics ---")
	var metrics []*models.Metrics
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.memStorage) == 0 {
		log.Println("Storage is empty.")
		return nil, errors.New("storage is empty")
	}

	for key, metric := range r.memStorage {
		metrics = append(metrics, metric)

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
	return metrics, nil
}

func (r *InMemoryRepo) Add(metrics []*models.Metrics) error {
	r.mu.Lock() // Блокируем мьютекс один раз для всей операции
	defer r.mu.Unlock()

	for _, metric := range metrics {
		oldmetric, ok := r.memStorage[metric.ID]
		if ok {
			switch metric.MType {
			case "gauge":
				if metric.Value != nil {
					oldmetric.Value = metric.Value
				}
			case "counter":
				if metric.Delta != nil {
					oldmetric.Delta = metric.Delta
				}
			}
		} else {
			// Если метрика новая, добавляем её в хранилище
			r.memStorage[metric.ID] = metric
		}
	}
	return nil
}

func (r *InMemoryRepo) Ping() error {
	return errors.New("cant ping postgres: current storage is inmemory")
}

func (r *InMemoryRepo) UpdateGauge(name string, value float64) (*models.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.memStorage[name]
	if !ok {
		return nil, errors.New("metric not found")
	}
	oldValue := r.memStorage[name].Value
	r.memStorage[name].Value = &value

	log.Printf("gauge %s updated from %v to %v\n", name, oldValue, r.memStorage[name].Value)

	return r.memStorage[name], nil
}

func (r *InMemoryRepo) UpdateCounter(name string, value int64) (*models.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.memStorage[name]
	if !ok {
		return nil, errors.New("metric not fount")
	}

	oldValue := *r.memStorage[name].Delta
	*r.memStorage[name].Delta += value

	log.Printf("counter %s updated from %v to %v\n", name, oldValue, *r.memStorage[name].Delta)

	return r.memStorage[name], nil
}

func (r *InMemoryRepo) Metric(name string) (*models.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metric, ok := r.memStorage[name]
	log.Println("name = ", name)
	if !ok {
		for _, val := range r.memStorage {
			log.Println(val)
		}

		return nil, errors.New("metric not found")
	}
	return metric, nil
}

func (r *InMemoryRepo) AddMetric(metric *models.Metrics) error {
	if metric == nil {
		return errors.New("metric is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	r.memStorage[metric.ID] = metric
	log.Println("Metric added to memStorage")
	return nil

}

func (r *InMemoryRepo) CheckExistence(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.memStorage[name]
	return ok
}
