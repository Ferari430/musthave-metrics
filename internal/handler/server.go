package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg"
)

type ServerHandler struct {
	Service *service.ServiceServer
}

type ServerHandlerDeps struct {
	Service *service.ServiceServer
}

func NewServerHandler(router *http.ServeMux, deps ServerHandlerDeps) {

	handler := &ServerHandler{
		Service: deps.Service,
	}

	router.HandleFunc("POST /update/{typeM}/{nameM}/{value}", handler.ProcessingMetric)

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

	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]

	log.Printf("metricType=%q, metricName=%q, metricValue=%q\n", metricType, metricName, metricValue)
	status, err := pkg.Validate(metricType, metricName, metricValue)
	if err != nil {
		pkg.ResponceHTTP(w, err.Error(), status)
		return
	}

	if err := handler.Service.AddMetrics(metricType, metricName, metricValue); err != nil {
		pkg.ResponceHTTP(w, "internal server error", http.StatusInternalServerError)
		return
	}

	pkg.ResponceHTTP(w, "ok", http.StatusOK)
}
