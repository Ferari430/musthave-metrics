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
	MetricJSON(metric *models.Metrics) bool                 //
}

type RepositoryAgent interface {
	Add(metrics models.MetricsAgent)
	Metrics() map[string]float64
}
