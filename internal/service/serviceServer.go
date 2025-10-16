package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/interfaces"
	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
)

type ServiceServer struct {
	ticker      *time.Ticker
	Repo        interfaces.Repository
	FileStorage interfaces.FileStorage
}

func NewServiceServer(ticker *time.Ticker, repo interfaces.Repository, fileStorage interfaces.FileStorage) *ServiceServer {
	return &ServiceServer{
		ticker:      ticker,
		Repo:        repo,
		FileStorage: fileStorage,
	}
}

func (s *ServiceServer) UpdateMetric(metric *models.Metrics) (*models.Metrics, error) {

	switch metric.MType {

	case models.Gauge:
		s.Repo.UpdateGauge(metric.ID, *metric.Value)
		return metric, nil
	case models.Counter:

		newmetric, err := s.Repo.UpdateCounter(metric.ID, *metric.Delta)
		if err != nil {

			return nil, err
		}
		return newmetric, nil
	default:
		return nil, errors.New("invalid metric type")
	}
}

func (s *ServiceServer) GetMetricJSON(metric *models.Metrics) (*models.Metrics, error) {
	if metric.Value != nil {
		fmt.Println("Value =", *metric.Value)
	} else {
		fmt.Println("Value is nil")
	}

	//get metric
	dbMetric, err := s.Repo.Metric(metric.ID)
	if err != nil {
		return nil, err
	}

	return dbMetric, nil
}

func (s *ServiceServer) AddMetric(metricType, metricName, metricValue string) (*models.Metrics, error) {

	switch metricType {
	case models.Gauge:
		ok := s.Repo.CheckExistence(metricName)
		if !ok {
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				return nil, err
			}

			metric := models.NewGauge(metricName, metricType, &value)
			err = s.Repo.AddMetric(metric)
			if err != nil {
				return nil, err
			}
			return metric, nil

		}
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return nil, err
		}
		existingMetric, err := s.Repo.UpdateGauge(metricName, value)
		if err != nil {
			return nil, err
		}
		return existingMetric, nil

	case models.Counter:
		ok := s.Repo.CheckExistence(metricName)
		if !ok {

			value, err := strconv.ParseInt(metricValue, 64, 10)
			if err != nil {
				return nil, err
			}
			metric := models.NewCounter(metricName, metricType, &value)
			err = s.Repo.AddMetric(metric)
			if err != nil {
				return nil, err
			}
			return metric, nil

		}
		value, err := strconv.ParseInt(metricValue, 64, 10)
		if err != nil {
			return nil, err
		}
		existingMetric, err := s.Repo.UpdateCounter(metricName, value)
		if err != nil {
			return nil, err
		}
		return existingMetric, nil

	default:
		return nil, errors.New("invalid metric type")
	}

}

func (s *ServiceServer) GetMetric(metricType, metricName string) (*models.Metrics, error) {
	statuscode, err := pkg.ValidateNameType(metricType, metricName)
	if err != nil {
		log.Println(statuscode)
		return nil, err
	}

	logger.Log.Info("cant [logger Zap")
	if statuscode != http.StatusOK {
		return nil, errors.New("invalid metric")
	}

	existingmetric, err := s.Repo.Metric(metricName)

	if err != nil {
		log.Println("cant find metric in storage")
		return nil, err
	}

	log.Println(existingmetric, err)

	return existingmetric, nil
}

func (s *ServiceServer) AllMetrics() []*models.HTMLMetricData {
	HTMLdata := make([]*models.HTMLMetricData, 0, 10)
	if s.Repo == nil {
		log.Println("Repo is nil")
		return HTMLdata
	}

	metrics, err := s.Repo.GetAll()
	if err != nil {
		log.Println(err)
		return nil
	}

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

// func (s *ServiceServer) AddJsonMetricsBatch(metrics []*models.Metrics) error {
// 	for _, metric := range metrics {

// 		switch metric.MType {
// 		case "gauge":
// 			log.Println("Now metric: ", *metric.Value)
// 		case "counter":
// 			log.Println("Now metric: ", *metric.Delta)
// 		}

// 		if metric.ID == "" {
// 			return errors.New("metric name is empty")
// 		}

// 	}

// 	s.Producer.WriteMetric(metrics)
// 	err := s.Repo.MetricJSONBatch(metrics)
// 	if err != nil {
// 		return errors.New("cant add metric in db")
// 	}

// 	return nil
// }

func (s *ServiceServer) AddJsonMetricsBatchTicker(metrics []*models.Metrics) error {

	select {
	case <-s.ticker.C:
		log.Println("tick")
		if err := s.FileStorage.Add(metrics); err != nil {
			return err
		}
		log.Println("Metrics written on file")
	default:
	}
	err := s.Repo.Add(metrics)
	if err != nil {
		return errors.New("cant add metric in db")
	}
	s.Repo.GetAll()
	return nil
}

func (s *ServiceServer) RestoreData() error {

	metrics, err := s.FileStorage.Restore()
	//можно вынести логику с ошибками в pkg
	if errors.Is(err, io.EOF) {
		log.Println("file store is empty")
	}
	if err != nil {
		return err
	}
	if metrics == nil {
		log.Println("no metrics to restore")
		return nil
	}

	err = s.Repo.Add(metrics)
	if err != nil {
		return err
	}
	return nil
}

// func (s *ServiceServer) PingPostgres() error {
// 	dsn := "postgres://postgres:pass@localhost:5432/postgres"
// 	err := s.Repo.Ping()
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("Postgres is alive")
// 	return nil
// }
