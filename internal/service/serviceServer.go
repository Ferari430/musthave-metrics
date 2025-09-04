package service

import (
	"errors"
	"log"
	"strconv"

	"github.com/Ferari430/musthave-metrics/internal/interfaces"
	models "github.com/Ferari430/musthave-metrics/internal/model"
)

// бизнес лоигка (сохранение в бд через интерфейс)
type ServiceServer struct {
	repo interfaces.Repository
}

func NewServiceServer(repo interfaces.Repository) *ServiceServer {
	return &ServiceServer{repo: repo}
}

func (s *ServiceServer) AddMetrics(metricType, metricName, metricValue string) error {

	switch metricType {
	//добавить к старому
	case models.Counter:
		intMetricValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("Cant parse metric value")
		}
		oldMetric, ok := s.repo.Get(metricName)
		if ok {
			*oldMetric.Delta += intMetricValue
			s.repo.Add(oldMetric)
		} else {
			metric := models.Metrics{
				ID:    metricName,
				MType: metricType,
				Delta: &intMetricValue,
			}
			s.repo.Add(&metric)

		}

	//обновить старое
	case models.Gauge:

		intMetricValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("Cant parse metric value")
		}

		oldMetric, ok := s.repo.Get(metricName)
		if ok {
			*oldMetric.Value = intMetricValue
			s.repo.Add(oldMetric)
		} else {
			metric := models.Metrics{
				ID:    metricName,
				MType: metricType,
				Value: &intMetricValue,
			}
			s.repo.Add(&metric)

		}

	default:
		log.Println("Unknown metric type")
		return errors.New("Unknown metric type")
	}

	return nil
}
