package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/interfaces"
	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"go.uber.org/zap"
)

type ServiceServer struct {
	ticker      *time.Ticker
	Repo        interfaces.Repository
	FileStorage interfaces.FileStorage
	logger      *logger.Logger
}

func NewServiceServer(ticker *time.Ticker, repo interfaces.Repository,
	fileStorage interfaces.FileStorage, logger *logger.Logger,
) *ServiceServer {
	return &ServiceServer{
		ticker:      ticker,
		Repo:        repo,
		FileStorage: fileStorage,
		logger:      logger,
	}
}

func (s *ServiceServer) UpdateMetric(metric *models.Metrics) (*models.Metrics, error) {
	op := "Service.UpdateMetric"

	s.logger.Debug("incoming metric",
		zap.String("operation", op),
		zap.Any("metric", metric),
	)

	switch metric.MType {
	case models.Gauge:
		s.Repo.UpdateGauge(metric.ID, *metric.Value)

		s.logger.Debug("gauge updated",
			zap.String("operation", op),
			zap.String("metric", metric.ID),
			zap.Float64("value", *metric.Value),
		)

		return metric, nil

	case models.Counter:
		newMetric, err := s.Repo.UpdateCounter(metric.ID, *metric.Delta)
		if err != nil {

			s.logger.Error("failed to update counter",
				zap.String("operation", op),
				zap.String("metric", metric.ID),
				zap.Error(err),
			)

			return nil, err
		}

		s.logger.Debug("counter updated",
			zap.String("operation", op),
			zap.String("metric", metric.ID),
			zap.Int64("delta", *metric.Delta),
		)

		return newMetric, nil

	default:
		s.logger.Warn("invalid metric type",
			zap.String("operation", op),
			zap.String("metric", metric.ID),
			zap.String("type", metric.MType),
		)
		return nil, errors.New("invalid metric type")
	}
}

func (s *ServiceServer) GetMetricJSON(metric *models.Metrics) (*models.Metrics, error) {
	op := "Service.GetMetricJSON"

	s.logger.Debug("get metric JSON",
		zap.String("operation", op),
		zap.String("metric", metric.ID),
		zap.String("type", metric.MType),
	)

	dbMetric, err := s.Repo.Metric(metric.ID)
	if err != nil {
		s.logger.Error("metric not found",
			zap.String("operation", op),
			zap.String("metric", metric.ID),
			zap.Error(err),
		)
		return nil, err
	}

	return dbMetric, nil
}

func (s *ServiceServer) AddMetric(metricType, metricName, metricValue string) (*models.Metrics, error) {
	op := "Service.AddMetric"

	s.logger.Debug("adding metric",
		zap.String("operation", op),
		zap.String("type", metricType),
		zap.String("name", metricName),
		zap.String("value", metricValue),
	)

	switch metricType {

	case models.Gauge:
		ok := s.Repo.CheckExistence(metricName)
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			s.logger.Error("failed to parse float gauge",
				zap.String("operation", op),
				zap.String("name", metricName),
				zap.String("value", metricValue),
				zap.Error(err),
			)
			return nil, err
		}

		if !ok {
			metric := models.NewGauge(metricName, metricType, &value)
			err = s.Repo.AddMetric(metric)
			if err != nil {
				s.logger.Error("failed to add gauge",
					zap.String("operation", op),
					zap.String("metric", metricName),
					zap.Error(err),
				)
				return nil, err
			}

			s.logger.Debug("gauge created",
				zap.String("operation", op),
				zap.String("metric", metricName),
				zap.Float64("value", value),
			)
			return metric, nil
		}

		existing, err := s.Repo.UpdateGauge(metricName, value)
		if err != nil {
			s.logger.Error("failed to update existing gauge",
				zap.String("operation", op),
				zap.String("metric", metricName),
				zap.Error(err),
			)
			return nil, err
		}
		return existing, nil

	case models.Counter:
		ok := s.Repo.CheckExistence(metricName)
		value, err := strconv.ParseInt(metricValue, 64, 10)
		if err != nil {
			s.logger.Error("failed to parse counter int",
				zap.String("operation", op),
				zap.String("name", metricName),
				zap.String("value", metricValue),
				zap.Error(err),
			)
			return nil, err
		}

		if !ok {
			metric := models.NewCounter(metricName, metricType, &value)
			err = s.Repo.AddMetric(metric)
			if err != nil {
				s.logger.Error("failed to add counter",
					zap.String("operation", op),
					zap.String("metric", metricName),
					zap.Error(err),
				)
				return nil, err
			}
			return metric, nil
		}

		existingMetric, err := s.Repo.UpdateCounter(metricName, value)
		if err != nil {
			s.logger.Error("failed to update counter",
				zap.String("operation", op),
				zap.String("metric", metricName),
				zap.Error(err),
			)
			return nil, err
		}
		return existingMetric, nil

	default:
		s.logger.Warn("invalid metric type",
			zap.String("operation", op),
			zap.String("metricType", metricType),
		)
		return nil, errors.New("invalid metric type")
	}
}

func (s *ServiceServer) GetMetric(metricType, metricName string) (*models.Metrics, error) {
	op := "Service.GetMetric"

	s.logger.Debug("get metric",
		zap.String("operation", op),
		zap.String("type", metricType),
		zap.String("name", metricName),
	)

	statuscode, err := pkg.ValidateNameType(metricType, metricName)
	if err != nil {
		s.logger.Warn("validation failed",
			zap.String("operation", op),
			zap.Int("statuscode", statuscode),
			zap.Error(err),
		)
		return nil, err
	}

	if statuscode != http.StatusOK {
		return nil, errors.New("invalid metric")
	}

	existingMetric, err := s.Repo.Metric(metricName)
	if err != nil {
		s.logger.Error("metric not found",
			zap.String("operation", op),
			zap.String("metric", metricName),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Debug("metric fetched successfully",
		zap.String("operation", op),
		zap.Any("metric", existingMetric),
	)

	return existingMetric, nil
}

func (s *ServiceServer) AllMetrics() []*models.HTMLMetricData {
	op := "Service.AllMetrics"

	HTMLdata := make([]*models.HTMLMetricData, 0, 10)
	if s.Repo == nil {

		s.logger.Error("repo is nil",
			zap.String("operation", op),
		)

		return HTMLdata
	}

	metrics, err := s.Repo.GetAll()
	if err != nil {
		s.logger.Error("failed to fetch metrics",
			zap.String("operation", op),
			zap.Error(err),
		)
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
				s.logger.Warn("counter value is nil",
					zap.String("operation", op),
					zap.String("metric", metric.ID),
				)
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
				s.logger.Warn("gauge value is nil",
					zap.String("operation", op),
					zap.String("metric", metric.ID),
				)
				continue
			}
			value := fmt.Sprintf("%v", *metric.Value)
			HTMLdata = append(HTMLdata, &models.HTMLMetricData{
				Name:  metric.ID,
				Value: value,
				Type:  metric.MType,
			})

		default:
			s.logger.Warn("unknown metric type",
				zap.String("operation", op),
				zap.String("metric", metric.ID),
			)
		}
	}

	for _, m := range HTMLdata {
		s.logger.Debug("html metric",
			zap.String("operation", op),
			zap.String("name", m.Name),
			zap.String("type", m.Type),
			zap.String("value", m.Value),
		)
	}

	return HTMLdata
}

func (s *ServiceServer) AddJsonMetricsBatchTicker(metrics []*models.Metrics) error {
	op := "Service.AddJsonMetricsBatchTicker"

	select {
	case <-s.ticker.C:
		s.logger.Debug("ticker tick",
			zap.String("operation", op),
		)

		if err := s.FileStorage.Add(metrics); err != nil {
			s.logger.Error("file write failed",
				zap.String("operation", op),
				zap.Error(err),
			)
			return err
		}

		s.logger.Debug("metrics written to file",
			zap.String("operation", op),
		)

	default:
	}

	err := s.Repo.Add(metrics)
	if err != nil {
		s.logger.Error("failed to add metrics to db",
			zap.String("operation", op),
			zap.Error(err),
		)
		return errors.New("cant add metric in db")
	}

	s.logger.Debug("db metrics added successfully",
		zap.String("operation", op),
	)

	s.Repo.GetAll()
	return nil
}

func (s *ServiceServer) RestoreData() error {
	op := "Service.RestoreData"

	metrics, err := s.FileStorage.Restore()

	if errors.Is(err, io.EOF) {
		s.logger.Info("file storage empty",
			zap.String("operation", op),
		)
	}

	if err != nil {
		s.logger.Error("restore failed",
			zap.String("operation", op),
			zap.Error(err),
		)
		return err
	}

	if metrics == nil {
		s.logger.Info("no metrics to restore",
			zap.String("operation", op),
		)
		return nil
	}

	s.logger.Debug("restored metrics from file",
		zap.String("operation", op),
		zap.Int("count", len(metrics)),
	)

	// Возможное восстановление в БД (в ТЗ отключено)
	return nil
}
