package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Ferari430/musthave-metrics/internal/interfaces"
	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg"
)

type ServiceServer struct {
	Repo interfaces.Repository
}

func NewServiceServer(Repo interfaces.Repository) *ServiceServer {
	return &ServiceServer{Repo: Repo}
}

func (s *ServiceServer) UpdateMetric(metric *models.Metrics) (*models.Metrics, error) {

	UpdatedMetric, err := s.Repo.Update(metric)
	if err != nil {

		return nil, err
	}

	return UpdatedMetric, nil
}

func (s *ServiceServer) GetMetricJSON(metric *models.Metrics) error {
	if metric.Value != nil {
		fmt.Println("Value =", *metric.Value)
	} else {
		fmt.Println("Value is nil")
	}

	ok := s.Repo.MetricJSON(metric)
	if !ok {
		//todo
		return errors.New("metric not found in database")
	}

	return nil

}

func (s *ServiceServer) AddMetrics(metricType, metricName, metricValue string) error {

	switch metricType {
	//добавить к старому
	case models.Counter:
		intMetricValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("Cant parse metric value")
		}
		oldMetric, ok := s.Repo.Metric(metricName)
		if ok {
			*oldMetric.Delta += intMetricValue
			s.Repo.Add(oldMetric)
		} else {
			metric := models.Metrics{
				ID:    metricName,
				MType: metricType,
				Delta: &intMetricValue,
			}
			s.Repo.Add(&metric)

		}

	//обновить старое
	case models.Gauge:

		intMetricValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("Cant parse metric value")
		}

		oldMetric, ok := s.Repo.Metric(metricName)
		if ok {
			*oldMetric.Value = intMetricValue
			s.Repo.Add(oldMetric)
		} else {
			metric := models.Metrics{
				ID:    metricName,
				MType: metricType,
				Value: &intMetricValue,
			}
			s.Repo.Add(&metric)

		}

	default:
		log.Println("Unknown metric type")
		return errors.New("Unknown metric type")
	}

	return nil
}

func (s *ServiceServer) GetMetric(metricType, metricName string) (*models.Metrics, error) {

	statuscode, err := pkg.ValidateNameType(metricType, metricName)
	if err != nil {
		return nil, err
	}

	if statuscode != http.StatusOK {
		return nil, errors.New("invalid metric")
	}

	existingmetric, ok := s.Repo.Metric(metricName)
	if !ok {
		return nil, errors.New("cant find metric in repo")
	}

	return existingmetric, nil
}

func (s *ServiceServer) AllMetrics() []*models.HTMLMetricData {
	HTMLdata := make([]*models.HTMLMetricData, 0, 10)
	if s.Repo == nil {
		log.Println("Repo is nil")
		return HTMLdata
	}

	metrics := s.Repo.Metrics()
	if metrics == nil {
		return HTMLdata
	}

	for _, metric := range metrics {
		if metric == nil {
			continue
		}

		switch metric.MType {
		case models.Counter:
			if metric.Delta == nil {
				log.Printf("counter %s is nil", metric.ID)
				continue
			}
			value := fmt.Sprintf("%v", *metric.Delta)
			HTMLdata = append(HTMLdata, &models.HTMLMetricData{
				Name:  metric.ID,
				Value: value,
				Type:  metric.MType,
			})
		case models.Gauge:
			if metric.Value == nil {
				log.Printf("gauge %s is nil", metric.ID)
				continue
			}
			value := fmt.Sprintf("%v", *metric.Value)
			HTMLdata = append(HTMLdata, &models.HTMLMetricData{
				Name:  metric.ID,
				Value: value,
				Type:  metric.MType,
			})
		default:
			log.Printf("[AllMetrics SERVICE] unknown type for metric %s", metric.ID)
		}
	}

	for _, m := range HTMLdata {
		log.Printf("<li>%s (%s) = %s</li>\n", m.Name, m.Type, m.Value)
	}

	return HTMLdata
}
