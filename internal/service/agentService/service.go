package agentService

import (
	"log"
	"runtime"
	"sync"
	"time"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	repositoryAgent "github.com/Ferari430/musthave-metrics/internal/repository/agent"
)

type AgentService struct {
	repo            *repositoryAgent.RepositoryAgent
	metricsChannel  chan []*models.Metrics
	pollCount       int64
	mu              sync.Mutex
	lastPollMetrics []*models.Metrics
	typedMetrics    map[string]*models.Metrics
}

func NewAgentService(repo *repositoryAgent.RepositoryAgent) *AgentService {

	channel := make(chan []*models.Metrics)
	agent := &AgentService{repo: repo,
		metricsChannel: channel,
		mu:             sync.Mutex{}}

	return agent
}

func (a *AgentService) MetricsChannel() chan []*models.Metrics {
	return a.metricsChannel
}

func float64Ptr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64       { return &v }

// сбор метрик и их типизация
func (a *AgentService) CollectMetrics(m *runtime.MemStats) []*models.Metrics {

	runtime.ReadMemStats(m)
	a.mu.Lock()
	a.pollCount++
	a.mu.Unlock()
	metrics := []*models.Metrics{
		//gauge
		{ID: "TotalAlloc", MType: "gauge", Value: float64Ptr(float64(m.TotalAlloc))},
		{ID: "HeapSys", MType: "gauge", Value: float64Ptr(float64(m.HeapSys))},

		//counter
		{ID: "Frees", MType: "counter", Delta: int64Ptr(int64(m.Frees))},
		{ID: "PollCounter", MType: "counter", Delta: int64Ptr(int64(a.pollCount))},
	}

	return metrics
}

// Переписать функцию с использованием каналов и не собирать метрики два раза.
func (a *AgentService) StartTicker(t1, t2 time.Ticker, m *runtime.MemStats, wg *sync.WaitGroup) {

	// Горутина для сбора метрик по интервалу t1 (pollInterval)
	go func() {
		for range t1.C {
			metrics := a.CollectMetrics(m)
			a.mu.Lock()
			a.lastPollMetrics = metrics // slice
			a.mu.Unlock()
			a.repo.Add(metrics) //map

		}
	}()

	// Горутина для отправки метрик по интервалу t2 (reportInterval)
	go func() {
		for range t2.C {
			a.mu.Lock()
			metricsToSend := a.lastPollMetrics
			a.mu.Unlock()
			if metricsToSend != nil {
				log.Println("send to server")
				a.metricsChannel <- metricsToSend
			} else {
				log.Println("No metrics collected yet to send.")
			}
		}
	}()
}
