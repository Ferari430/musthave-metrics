package interfaces

import (
	models "github.com/Ferari430/musthave-metrics/internal/model"
)

type Repository interface {
	Add(metrics *models.Metrics) //
	Metric(name string) (*models.Metrics, bool)
	Metrics() map[string]*models.Metrics                    //
	PrintAll()                                              //
	Update(metric *models.Metrics) (*models.Metrics, error) //
	MetricJSON(metric *models.Metrics) error                //
	MetricJSONBatch(metrics []*models.Metrics) error
}

type RepositoryAgent interface {
	Add(metrics models.MetricsAgent)
	Metrics() map[string]float64
}

type Postgres interface {
	Ping(dsn string) error
}

type GeneralRepo interface {
	Add(metrics *models.Metrics)
}
