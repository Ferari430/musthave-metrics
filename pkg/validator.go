package pkg

import (
	"errors"
	"net/http"
	"strconv"

	models "github.com/Ferari430/musthave-metrics/internal/model"
)

var (
	ErrNoMetricName = errors.New("metric name is required")
	ErrUnknownType  = errors.New("unknown metric type")
	ErrInvalidValue = errors.New("invalid metric value")
)

func Validate(metricType, metricName, metricValue string) (int, error) {
	if metricName == "" {
		return http.StatusNotFound, errors.New("metric name is required")
	}

	if metricType != models.Counter && metricType != models.Gauge {
		return http.StatusBadRequest, errors.New("unknown metric type")
	}

	if metricValue == "" {
		return http.StatusBadRequest, errors.New("metric value is required")
	}

	switch metricType {
	case models.Counter:
		if _, err := strconv.ParseInt(metricValue, 10, 64); err != nil {
			return http.StatusBadRequest, errors.New("value for counter must be int")
		}
	case models.Gauge:
		if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
			return http.StatusBadRequest, errors.New("value for gauge must be float")
		}
	}

	return http.StatusOK, nil
}
