package app

import (
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	repositoryAgent "github.com/Ferari430/musthave-metrics/internal/repository/agent"
	"github.com/Ferari430/musthave-metrics/internal/service/agentService"
	"github.com/Ferari430/musthave-metrics/pkg"
	gopsutilMem "github.com/shirou/gopsutil/v3/mem"
)

func StartApp(portServer, key string, pollInterval, reportInterval int, hashingFlag bool, rate_limit int) {
	db := repositoryAgent.NewInMemoryAgentDB()
	client := &http.Client{Timeout: 10 * time.Second}

	memStat, err := gopsutilMem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}

	sema := pkg.NewSemaphore(rate_limit)

	agentService := agentService.NewAgentService(db, memStat, sema)
	t1 := time.NewTicker(time.Second * time.Duration(pollInterval))
	t2 := time.NewTicker(time.Second * time.Duration(reportInterval))
	m := runtime.MemStats{}
	wg := sync.WaitGroup{}

	wg.Add(2)
	go agentService.StartAgent(*t1, *t2, &m, &wg)
	sender := handler.NewAgentSender(agentService, client, portServer, hashingFlag, key) // handler layer
	go sender.Consumer(&wg)
	wg.Wait()

	defer t1.Stop()
	defer t2.Stop()

}
