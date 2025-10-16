package pkg

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func ConfigurateServer() (string, int, string, bool, string, bool) {

	var (
		portServer        string
		store_interval    int
		file_storage_path string
		restore           bool
		dsn               string
		fileStorage       bool
	)

	// postgres  file  inmemory
	flag.StringVar(&portServer, "port", ":8080", "port for server")
	flag.IntVar(&store_interval, "i", 20, "interval for saving metrics in file")
	flag.StringVar(&file_storage_path, "f", "fileStorage.json", "file name storage")
	flag.BoolVar(&restore, "r", true, "enable restore")
	flag.StringVar(&dsn, "d", "postgres://postgres:postgres@localhost:5432/postgres", "dsn for connecting to postgres")
	flag.BoolVar(&fileStorage, "fs", true, "enable fileStorage")

	flag.Parse()

	if envAddres := os.Getenv("PORT"); envAddres != "" {
		portServer = envAddres
	}

	stringStore_interval, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		intVal, _ := strconv.Atoi(stringStore_interval)
		store_interval = intVal
	} else {
		log.Println("env var STORE_INTERVAL not found, using flag or default value:", store_interval)
	}

	sPath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		file_storage_path = sPath
	} else {
		log.Println("env var FILE_STORAGE_PATH not found, using flag or default value:", file_storage_path)
	}

	r, ok := os.LookupEnv("RESTORE")
	if ok {
		boolVal, _ := strconv.ParseBool(r)
		restore = boolVal
	} else {
		log.Println("env var RESTORE not found, using flag or default value:", restore)
	}

	d, ok := os.LookupEnv("DSN")
	if ok {
		dsn = d
	} else {
		log.Println("env var DSN not found, using flag or default value:", "localhost:5432")
	}

	fs, ok := os.LookupEnv("FILE_STORAGE")
	if ok {
		boolVal, _ := strconv.ParseBool(fs)
		fileStorage = boolVal
	} else {
		log.Println("env var FILE_STORAGE not found, using flag or default value:", fileStorage)
	}

	return portServer, store_interval, file_storage_path, restore, dsn, fileStorage
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
