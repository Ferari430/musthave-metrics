package interfaces

import (
	models "github.com/Ferari430/musthave-metrics/internal/model"
)

type FileStorage interface {
	Add(metrics []*models.Metrics) error
	GetAll() ([]*models.Metrics, error)
	Ping(fpath string) error
	Restore() ([]*models.Metrics, error)
}

type Repository interface {
	Add(metrics []*models.Metrics) error                             // добавить список метрик
	GetAll() ([]*models.Metrics, error)                              // получить все метрики
	Ping() error                                                     // пинг
	Metric(name string) (*models.Metrics, error)                     // получить метрику по имени
	AddMetric(metric *models.Metrics) error                          // добавить метрику по имени
	UpdateCounter(name string, value int64) (*models.Metrics, error) // апдейт метрики
	UpdateGauge(name string, value float64) (*models.Metrics, error) // апдейт метрики
	CheckExistence(name string) bool                                 // проверить существование метрики
}
