package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Ferari430/musthave-metrics/internal/service/agentService"
)

type AgentSender struct {
	client  *http.Client
	service *agentService.AgentService
}

func NewAgentSender(service *agentService.AgentService, client *http.Client) *AgentSender {
	return &AgentSender{client: client, service: service}
}

func (a *AgentSender) Consumer(wg *sync.WaitGroup) {
	// Эта горутина работает в течение всей жизни приложения.
	// wg.Done() не вызывается, что позволяет wg.Wait() в main блокировать программу бессрочно.
	channel := a.service.MetricsChannel()
	for metric := range channel {
		log.Printf("consumer: sending http request to server %v", metric)
		a.SendHTTP(metric)
	}
}

func (a *AgentSender) SendHTTP(metrics map[string]float64) {
	for key, val := range metrics {
		var url string

		metricType := "gauge"
		url = fmt.Sprintf("http://localhost:8080/update/%v/%v/%v", metricType, key, val)
		log.Println("sending request to " + url)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "text/plain")
		resp, err := a.client.Do(req)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		resp.Body.Close()
		log.Println("response Status:", resp.Status)
	}
}
