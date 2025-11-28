package main

import (
	"github.com/Ferari430/musthave-metrics/cmd/agent/app"
	"github.com/Ferari430/musthave-metrics/pkg"
)

func main() {
	portServer, pollInterval, reportInterval, hashingFlag, key, rate_limit := pkg.ConfigurateAgent()
	app.StartApp(portServer, key, pollInterval, reportInterval, hashingFlag, rate_limit)
}
