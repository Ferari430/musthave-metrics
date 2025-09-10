package pkg

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func ConfigurateServer() string {

	var (
		portServer string
	)

	flag.StringVar(&portServer, "port", ":8080", "port for server")

	flag.Parse()

	if envAddres := os.Getenv("PORT"); envAddres != "" {
		portServer = envAddres
	}

	return portServer
}

func ConfigurateAgent() (string, int64, int64) {

	var (
		portServer     string
		pollInterval   int64
		reportInterval int64
	)

	flag.Int64Var(&pollInterval, "pollInterval", 3, "PollIntervalValue in sec")
	flag.Int64Var(&reportInterval, "reportInterval", 5, "ReportIntervalValue in sec")
	flag.StringVar(&portServer, "port", "8080", "port for server")

	flag.Parse()

	strPollInterval := os.Getenv("pollInterval")
	intPollinterval, err := strconv.ParseInt(strPollInterval, 10, 64)

	if err == nil {
		pollInterval = intPollinterval
	}

	strReportInterval := os.Getenv("pollInterval")
	intReportInterval, err := strconv.ParseInt(strReportInterval, 10, 64)

	if err == nil {
		reportInterval = intReportInterval
	}

	if envAddres := os.Getenv("PORT"); envAddres != "" {
		portServer = envAddres
	}

	log.Printf("pollInterval: %v", pollInterval)
	log.Printf("reportInterval: %v", reportInterval)
	log.Printf("portServer: %v", portServer)
	return portServer, pollInterval, reportInterval
}
