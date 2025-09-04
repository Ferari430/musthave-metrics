package agentService

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	repositoryAgent "github.com/Ferari430/musthave-metrics/internal/repository/agent"
)

type AgentService struct {
	repo            *repositoryAgent.RepositoryAgent
	metricsChannel  chan map[string]float64
	pollCount       int64
	mu              sync.Mutex
	lastPollMetrics map[string]float64
}

func NewAgentService(repo *repositoryAgent.RepositoryAgent) *AgentService {

	channel := make(chan map[string]float64)
	agent := &AgentService{repo: repo,
		metricsChannel: channel,
		mu:             sync.Mutex{}}

	return agent
}

func (a *AgentService) MetricsChannel() chan map[string]float64 {
	return a.metricsChannel
}

func (a *AgentService) CollectMetrics(m *runtime.MemStats) map[string]float64 {
	runtime.ReadMemStats(m)
	a.mu.Lock()
	a.pollCount++
	a.mu.Unlock()
	metrics := map[string]float64{
		"Alloc":         float64(m.Alloc),
		"TotalAlloc":    float64(m.TotalAlloc),
		"Sys":           float64(m.Sys),
		"Frees":         float64(m.Frees),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapSys":       float64(m.HeapSys),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"GCSys":         float64(m.GCSys),
		"NumGC":         float64(m.NumGC),
		"GCCPUFraction": m.GCCPUFraction,
		"RandomValue":   rand.Float64(),
		"PollCount":     float64(a.pollCount),
	}

	log.Println("-------------------Metrics collected-------------------:")
	log.Println("-------------------------------", metrics["PollCount"], "---------------------------------")
	log.Println(metrics)
	log.Println("---------------------------------------------------------")
	return metrics
}

// Переписать функцию с использованием каналов и не собирать метрики два раза.
func (a *AgentService) StartTicker(t1, t2 time.Ticker, m *runtime.MemStats, wg *sync.WaitGroup) {

	// Горутина для сбора метрик по интервалу t1 (pollInterval)
	go func() {
		for range t1.C {
			metrics := a.CollectMetrics(m)
			a.mu.Lock()
			a.lastPollMetrics = metrics
			a.mu.Unlock()
			a.repo.Add(metrics)
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
