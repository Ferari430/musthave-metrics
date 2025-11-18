package inMemoryStorage

import (
	"errors"
	"fmt"
	"log"
	"sync"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"go.uber.org/zap"
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
	op := "InMemoryRepo.GetAll"

	log.Println("--- Current state of metrics ---")
	var metrics []*models.Metrics
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.memStorage) == 0 {
		r.logger.Debug("inMemory storage is empty", zap.String("operation", op), zap.Error(errors.New("storage is empty")))
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
			r.memStorage[metric.ID] = metric
		}
	}
	return nil
}

func (r *InMemoryRepo) Ping() error {
	return errors.New("cant ping postgres: current storage is inmemory")
}

func (r *InMemoryRepo) UpdateGauge(name string, value float64) (*models.Metrics, error) {
	op := "InMemoryRepo.UpdateGauge"

	r.mu.Lock()
	defer r.mu.Unlock()

	metric, ok := r.memStorage[name]
	if !ok {
		r.logger.Warn("gauge not found",
			zap.String("operation", op),
			zap.String("metric", name),
		)
		return nil, errors.New("metric not found")
	}

	oldValue := metric.Value
	metric.Value = &value

	r.logger.Debug("gauge updated",
		zap.String("operation", op),
		zap.String("metric", name),
		zap.Any("old_value", oldValue),
		zap.Float64("new_value", value),
	)

	return metric, nil
}

func (r *InMemoryRepo) UpdateCounter(name string, value int64) (*models.Metrics, error) {
	op := "InMemoryRepo.UpdateCounter"

	r.mu.Lock()
	defer r.mu.Unlock()

	metric, ok := r.memStorage[name]
	if !ok {
		r.logger.Warn("counter not found",
			zap.String("operation", op),
			zap.String("metric", name),
		)
		return nil, errors.New("metric not fount")
	}

	oldValue := *metric.Delta
	*metric.Delta += value

	r.logger.Debug("counter updated",
		zap.String("operation", op),
		zap.String("metric", name),
		zap.Int64("old_value", oldValue),
		zap.Int64("new_value", *metric.Delta),
	)

	return metric, nil
}

func (r *InMemoryRepo) Metric(name string) (*models.Metrics, error) {
	op := "InMemoryRepo.Metric"

	r.mu.RLock()
	defer r.mu.RUnlock()

	metric, ok := r.memStorage[name]

	if !ok {
		r.logger.Warn("metric not found",
			zap.String("operation", op),
			zap.String("metric", name),
			zap.Int("total_metrics", len(r.memStorage)),
		)
		return nil, errors.New("metric not found")
	}

	r.logger.Debug("metric fetched",
		zap.String("operation", op),
		zap.String("metric", name),
		zap.Any("value", metric),
	)

	return metric, nil
}

func (r *InMemoryRepo) AddMetric(metric *models.Metrics) error {
	op := "InMemoryRepo.AddMetric"

	if metric == nil {
		r.logger.Error("nil metric passed",
			zap.String("operation", op),
		)
		return errors.New("metric is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.memStorage[metric.ID] = metric

	r.logger.Info("metric added",
		zap.String("operation", op),
		zap.String("metric", metric.ID),
		zap.String("type", metric.MType),
	)

	return nil
}

func (r *InMemoryRepo) CheckExistence(name string) bool {
	op := "InMemoryRepo.CheckExistence"

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.memStorage[name]

	r.logger.Debug("checked existence",
		zap.String("operation", op),
		zap.String("metric", name),
		zap.Bool("exists", ok),
	)

	return ok
}
