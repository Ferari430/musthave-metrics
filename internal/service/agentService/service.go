package agentService

import (
	"log"
	"runtime"
	"sync"
	"time"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	repositoryAgent "github.com/Ferari430/musthave-metrics/internal/repository/agent"
	"github.com/Ferari430/musthave-metrics/pkg"
	cpu "github.com/shirou/gopsutil/v3/cpu"
	gopsutilMem "github.com/shirou/gopsutil/v3/mem"
)

type AgentService struct {
	repo            *repositoryAgent.RepositoryAgent
	metricsChannel  chan []*models.Metrics
	pollCount       int64
	mu              sync.Mutex
	lastPollMetrics []*models.Metrics
	typedMetrics    map[string]*models.Metrics
	memStat         *gopsutilMem.VirtualMemoryStat
	sema            *pkg.Semaphore
}

func NewAgentService(repo *repositoryAgent.RepositoryAgent, mem *gopsutilMem.VirtualMemoryStat, sema *pkg.Semaphore) *AgentService {

	channel := make(chan []*models.Metrics)
	agent := &AgentService{repo: repo,
		metricsChannel: channel,
		mu:             sync.Mutex{},
		memStat:        mem,
		sema:           sema,
	}

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

	cpuPersentage, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println(err)
	}

	totalMemory := models.NewGauge("TotalMemory", "gauge", float64Ptr(float64(a.memStat.Total)))
	freeMemory := models.NewGauge("FreeMemory", "gauge", float64Ptr(float64(a.memStat.Free)))
	usedMemoryPersentage := models.NewGauge("UsedMemoryPercentage", "gauge", &a.memStat.UsedPercent)
	usedCpuPersentage := models.NewGauge("UsedCpuPercentage", "gauge", &cpuPersentage[0])
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
		{ID: totalMemory.ID, MType: totalMemory.MType, Value: totalMemory.Value},
		{ID: freeMemory.ID, MType: freeMemory.MType, Value: freeMemory.Value},
		{ID: usedMemoryPersentage.ID, MType: usedMemoryPersentage.MType, Value: usedMemoryPersentage.Value},
		{ID: usedCpuPersentage.ID, MType: usedCpuPersentage.MType, Value: usedCpuPersentage.Value},
	}

	return metrics
}

func (a *AgentService) StartAgent(t1, t2 time.Ticker, m *runtime.MemStats, wg *sync.WaitGroup) {

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

	// Горутины для отправки метрик по интервалу t2 (reportInterval)
	for range t2.C {
		go func() {
			a.sema.Acquire()
			defer a.sema.Release()
			a.mu.Lock()
			metricsToSend := a.lastPollMetrics
			a.mu.Unlock()
			if metricsToSend != nil {
				log.Println("send to server")
				a.metricsChannel <- metricsToSend
			} else {
				log.Println("No metrics collected yet to send.")
			}
		}()
	}

}
