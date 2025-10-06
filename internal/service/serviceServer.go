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
	Repo     interfaces.Repository
	Producer *pkg.Producer
	ticker   *time.Ticker
	Consumer *pkg.Consumer
	Postgres interfaces.Postgres
}

func NewServiceServer(Repo interfaces.Repository, producer *pkg.Producer,
	ticker *time.Ticker, consumer *pkg.Consumer, postgres interfaces.Postgres) *ServiceServer {
	return &ServiceServer{Repo: Repo,
		Producer: producer,
		ticker:   ticker,
		Consumer: consumer,
		Postgres: postgres,
	}
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

	err := s.Repo.MetricJSON(metric)
	if err != nil {
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

	logger.Log.Info("cant [logger Zap")

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

func (s *ServiceServer) AddJsonMetricsBatch(metrics []*models.Metrics) error {
	for _, metric := range metrics {

		switch metric.MType {
		case "gauge":
			log.Println("Now metric: ", *metric.Value)
		case "counter":
			log.Println("Now metric: ", *metric.Delta)
		}

		if metric.ID == "" {
			return errors.New("metric name is empty")
		}

	}

	s.Producer.WriteMetric(metrics)
	err := s.Repo.MetricJSONBatch(metrics)
	if err != nil {
		return errors.New("cant add metric in db")
	}

	return nil
}

func (s *ServiceServer) AddJsonMetricsBatchTicker(metrics []*models.Metrics) error {

	select {
	case <-s.ticker.C:
		if err := s.Producer.WriteMetric(metrics); err != nil {
			return err
		}
		log.Println("Metrics written on file")
	default:
	}
	err := s.Repo.MetricJSONBatch(metrics)
	if err != nil {
		return errors.New("cant add metric in db")
	}

	return nil
}

func (s *ServiceServer) RestoreData() error {
	metrics, err := s.Consumer.Restore()
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

	err = s.Repo.MetricJSONBatch(metrics)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceServer) PingPostgres() error {
	dsn := "postgres://postgres:pass@localhost:5432/postgres"
	err := s.Postgres.Ping(dsn)
	if err != nil {
		return err
	}
	log.Println("Postgres is alive")
	return nil
}
