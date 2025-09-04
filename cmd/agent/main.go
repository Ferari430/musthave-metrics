package main

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	repositoryAgent "github.com/Ferari430/musthave-metrics/internal/repository/agent"
	"github.com/Ferari430/musthave-metrics/internal/service/agentService"
)

func main() {

	db := repositoryAgent.NewInMemoryAgentDB()
	client := &http.Client{Timeout: 5 * time.Second}
	// service := agentService.NewAgentService(db)
	// sender := handler.NewAgentSender(service, client)
	agentService := agentService.NewAgentService(db)
	t1 := time.NewTicker(time.Second * 3)
	t2 := time.NewTicker(time.Second * 5)
	m := runtime.MemStats{}
	wg := sync.WaitGroup{}

	wg.Add(2)
	go agentService.StartTicker(*t1, *t2, &m, &wg)
	sender := handler.NewAgentSender(agentService, client)
	go sender.Consumer(&wg)
	wg.Wait()

	// pollInterval := 2 * time.Second
	// reportInterval := 10 * time.Second
	// collectTicker := time.NewTicker(pollInterval)
	// sendTicker := time.NewTicker(reportInterval)
	// defer collectTicker.Stop()
	// defer sendTicker.Stop()
	defer t1.Stop()
	defer t2.Stop()
}
