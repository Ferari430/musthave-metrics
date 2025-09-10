package main

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	repositoryAgent "github.com/Ferari430/musthave-metrics/internal/repository/agent"
	"github.com/Ferari430/musthave-metrics/internal/service/agentService"
	"github.com/Ferari430/musthave-metrics/pkg"
)

func main() {
	portServer, pollInterval, reportInterval := pkg.ConfigurateAgent()
	db := repositoryAgent.NewInMemoryAgentDB()
	client := &http.Client{Timeout: 5 * time.Second}
	// service := agentService.NewAgentService(db)
	// sender := handler.NewAgentSender(service, client)
	agentService := agentService.NewAgentService(db)
	t1 := time.NewTicker(time.Second * time.Duration(pollInterval))
	t2 := time.NewTicker(time.Second * time.Duration(reportInterval))
	m := runtime.MemStats{}
	wg := sync.WaitGroup{}

	wg.Add(2)
	go agentService.StartTicker(*t1, *t2, &m, &wg)
	sender := handler.NewAgentSender(agentService, client, portServer)
	go sender.Consumer(&wg)
	wg.Wait()

	defer t1.Stop()
	defer t2.Stop()
}
