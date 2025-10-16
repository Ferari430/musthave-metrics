package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type ServerHandler struct {
	Service *service.ServiceServer
}

type ServerHandlerDeps struct {
	Service *service.ServiceServer
}

func NewServerHandler(router *chi.Mux, deps ServerHandlerDeps) {

	handler := &ServerHandler{
		Service: deps.Service,
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

	log.Printf("metric from agent: %v\n", metric)
	var updatedMetric *models.Metrics
	updatedMetric, err = handler.Service.UpdateMetric(metric)

	log.Println("updated metric:", updatedMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// responce
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedMetric)
	if err != nil {
		log.Println(err)
		return
	}
}

// get metric json by name and type TEST THIS
func (handler *ServerHandler) Value(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("metric from agent: %v\n", metric)
	dbMetric, err := handler.Service.GetMetricJSON(&metric)
	if err != nil {
		log.Println("cant get metric from db:", err)
		return
	}

	// responce
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dbMetric)
	if err != nil {
		log.Println(err)
		return
	}
}

// handler for agent
func (handler *ServerHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := strings.ToLower(chi.URLParam(r, "typeM"))
	metricName := chi.URLParam(r, "nameM")

	log.Printf("metricType=%q, metricName=%q\n", metricType, metricName)

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

func (handler *ServerHandler) ProcessingMetric(w http.ResponseWriter, r *http.Request) {

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
	log.Printf("metricType=%q, metricName=%q, metricValue=%q\n", metricType, metricName, metricValue)
	status, err := pkg.Validate(metricType, metricName, metricValue)
	if err != nil {
		pkg.ResponceHTTP(w, err.Error(), status)
		return
	}

	metric, err := handler.Service.AddMetric(metricType, metricName, metricValue)
	if err != nil {
		pkg.ResponceHTTP(w, "internal server error", http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("metric add to db: %v", metric)
	pkg.ResponceHTTP(w, message, http.StatusOK)
}

func (handler *ServerHandler) MetricsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

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
	log.Println("JSONCompressedHandler")

	if r.Method != http.MethodPost {

		http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
	}

	if r.Header.Get("Content-Type") != "application/json" {

		http.Error(w, "content-type must be json", http.StatusBadRequest)
	}

	var metrics []*models.Metrics
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		log.Println(err)
		log.Println("cant decode body")
		http.Error(w, "cant decode body", http.StatusBadRequest)
		return
	}

	for _, val := range metrics {
		if val.MType == "gauge" {

			log.Println(val.ID, val.MType, *val.Value)
		}
	}

	err = handler.Service.AddJsonMetricsBatchTicker(metrics)
	if err != nil {
		log.Println(err)
		return
	}

}

func (handler *ServerHandler) Ping(w http.ResponseWriter, r *http.Request) {

	err := handler.Service.Repo.Ping()

	if err != nil {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "postgres connected")

}
