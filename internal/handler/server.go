package handler

import (
	"bytes"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/hash"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ServerHandler struct {
	Service *service.ServiceServer
	Config  HandlerConfig
	Logger  *logger.Logger
}

type HandlerConfig struct {
	Hashing bool
	Key     string
}

type ServerHandlerDeps struct {
	Service *service.ServiceServer
	Config  HandlerConfig
	Logger  *logger.Logger
}

func NewServerHandler(router *chi.Mux, deps ServerHandlerDeps, logger *logger.Logger) {

	handler := &ServerHandler{
		Service: deps.Service,
		Config:  deps.Config,
		Logger:  logger,
	}
	router.Post("/update/{typeM}/{nameM}/{value}", logger.RequestLogger(handler.ProcessingMetric))
	router.Get("/update/{typeM}/{nameM}", logger.RequestLogger(handler.GetMetric))
	router.Get("/", logger.RequestLogger(handler.MetricsPage))
	router.Post("/update", logger.RequestLogger(handler.Update))
	router.Post("/value", logger.RequestLogger(handler.Value))
	router.Post("/valueJ", logger.RequestLogger(pkg.GzipMiddleware(handler.JSONCompressedMetric)))
	router.Get("/ping", logger.RequestLogger(handler.Ping))

}

func (handler *ServerHandler) Update(w http.ResponseWriter, r *http.Request) {
	op := "Server.UpdateHandler"

	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-Type")

	if !strings.HasPrefix(ct, "application/json") {

		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var metric *models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, "cant decode body", http.StatusBadRequest)
		return
	}

	handler.Logger.Debug("metric from agent",
		zap.String("operation", op),
		zap.Any("metric", metric),
	)

	var updatedMetric *models.Metrics
	updatedMetric, err = handler.Service.UpdateMetric(metric)

	handler.Logger.Debug("updated metric",
		zap.String("operation", op),
		zap.Any("updatedMetric", updatedMetric),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// responce
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedMetric)
	if err != nil {
		handler.Logger.Debug("cant encode updated metric for responce",
			zap.String("operation", op),
			zap.Error(err),
		)
		return
	}
}

// get metric json by name and type TEST THIS
func (handler *ServerHandler) Value(w http.ResponseWriter, r *http.Request) {
	op := "Server.ValueHandler"

	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
		return
	}

	ct := r.Header.Get("Content-Type")

	if !strings.HasPrefix(ct, "application/json") {

		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)

	if err != nil {
		http.Error(w, "cant decode body", http.StatusBadRequest)
		return
	}

	handler.Logger.Debug("got metric from agent",
		zap.String("operation", op),
		zap.Any("metric", metric),
	)

	dbMetric, err := handler.Service.GetMetricJSON(&metric)
	if err != nil {
		handler.Logger.Error("cant get metric from db",
			zap.String("operation", op),
			zap.Error(err),
		)
		return
	}

	// responce
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dbMetric)
	if err != nil {
		handler.Logger.Debug("cant encode metric",
			zap.String("operation", op),
			zap.Error(err),
		)
		return
	}
}

// handler for agent
func (handler *ServerHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	op := "Server.GetMetricHandler"

	defer r.Body.Close()
	metricType := strings.ToLower(chi.URLParam(r, "typeM"))
	metricName := chi.URLParam(r, "nameM")

	handler.Logger.Debug("getting metric",
		zap.String("operation", op),
		zap.String("metricType", metricType),
		zap.String("metricName", metricName),
	)

	existingmetric, err := handler.Service.GetMetric(metricType, metricName)
	if err != nil {

		http.Error(w, "invalid metric", http.StatusBadRequest)
		return
	}
	var txt string
	switch metricType {
	case models.Counter:
		txt = fmt.Sprintf("%v", *existingmetric.Delta)
	case models.Gauge:
		txt = fmt.Sprintf("%v", *existingmetric.Value)
	default:
		txt = fmt.Sprintf("unknown metric type: %q", metricType)
		http.Error(w, txt, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, txt)

}

func (handler *ServerHandler) Metric(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

func (handler *ServerHandler) ProcessingMetric(w http.ResponseWriter, r *http.Request) {
	op := "Server.ProcessingMetricHandler"

	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "text/plain" {
		pkg.ResponceHTTP(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		pkg.ResponceHTTP(w, "not enough params", http.StatusNotFound)
		return
	}

	metricType := strings.ToLower(chi.URLParam(r, "typeM"))
	metricName := strings.ToLower(chi.URLParam(r, "nameM"))
	metricValue := strings.ToLower(chi.URLParam(r, "value"))

	handler.Logger.Debug("processing metric from agent",
		zap.String("operation", op),
		zap.String("metricType", metricType),
		zap.String("metricName", metricName),
		zap.String("metricValue", metricValue),
	)

	status, err := pkg.Validate(metricType, metricName, metricValue)
	if err != nil {
		pkg.ResponceHTTP(w, err.Error(), status)
		return
	}

	metric, err := handler.Service.AddMetric(metricType, metricName, metricValue)
	if err != nil {
		handler.Logger.Error("failed to add metric",
			zap.String("operation", op),
			zap.Error(err),
		)
		pkg.ResponceHTTP(w, "internal server error", http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("metric add to db: %v", metric)
	pkg.ResponceHTTP(w, message, http.StatusOK)
}

func (handler *ServerHandler) MetricsPage(w http.ResponseWriter, r *http.Request) {
	op := "Server.MetricsPageHandler"

	defer r.Body.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	handler.Logger.Debug("rendering metrics page",
		zap.String("operation", op),
	)

	metrics := handler.Service.AllMetrics() // безопасно возвращает срез HTMLMetricData

	// Начало HTML
	fmt.Fprint(w, "<!DOCTYPE html>\n<html>\n<head>\n<title>Metrics</title>\n</head>\n<body>\n")
	fmt.Fprint(w, "<h1>Current Metrics</h1>\n<ul>\n")

	for _, m := range metrics {
		name := html.EscapeString(m.Name)
		value := html.EscapeString(m.Value)
		typ := html.EscapeString(m.Type)
		fmt.Fprintf(w, "<li>%s (%s) = %s</li>\n", name, typ, value)
	}

	fmt.Fprint(w, "</ul>\n</body>\n</html>")
}

func (handler *ServerHandler) JSONCompressedMetric(w http.ResponseWriter, r *http.Request) {
	op := "Server.JSONCompressedMetricHandler"

	defer r.Body.Close()

	handler.Logger.Debug("JSONCompressedHandler",
		zap.String("operation", op),
	)

	if r.Method != http.MethodPost {

		http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {

		http.Error(w, "content-type must be json", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		handler.Logger.Error("failed to read request body",
			zap.String("operation", op),
			zap.Error(err),
		)
		return
	}

	if handler.Config.Hashing {
		expectedHash := hash.ComputeHash(data, handler.Config.Key)
		recievedHash := r.Header.Get("HashSHA256")

		if recievedHash != "" {
			if !hmac.Equal([]byte(expectedHash), []byte(recievedHash)) {
				handler.Logger.Warn("hash mismatch",
					zap.String("operation", op),
					zap.String("expectedHash", expectedHash),
					zap.String("receivedHash", recievedHash),
				)
				http.Error(w, "hash not equal", http.StatusBadRequest)
				return
			}
			handler.Logger.Debug("signature true",
				zap.String("operation", op),
			)
		}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))
	var metrics []*models.Metrics
	err = json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		handler.Logger.Error("cant decode body",
			zap.String("operation", op),
			zap.Error(err),
		)
		http.Error(w, "cant decode body", http.StatusBadRequest)
		return
	}

	err = handler.Service.AddJsonMetricsBatchTicker(metrics)
	if err != nil {
		handler.Logger.Error("cant add metrics batch",
			zap.String("operation", op),
			zap.Error(err),
		)
		return
	}
}

func (handler *ServerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	op := "Server.PingHandler"

	defer r.Body.Close()
	err := handler.Service.Repo.Ping()

	if err != nil {
		handler.Logger.Error("postgres ping error",
			zap.String("operation", op),
			zap.Error(err),
		)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, err.Error())
		return
	}

	handler.Logger.Debug("postgres connected",
		zap.String("operation", op),
	)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "postgres connected")
}
