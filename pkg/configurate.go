package pkg

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func ConfigurateServer() (string, int, string, bool, string, bool, bool, string) {
	var (
		portServer      string
		storeInterval   int
		fileStoragePath string
		restore         bool
		dsn             string
		fileStorage     bool
		enableHashing   bool
		key             string
	)

	flag.StringVar(&portServer, "port", ":8080", "port for server")
	flag.IntVar(&storeInterval, "i", 20, "interval for saving metrics in file")
	flag.StringVar(&fileStoragePath, "f", "fileStorage.json", "file name storage")
	flag.BoolVar(&restore, "r", true, "enable restore")
	flag.StringVar(&dsn, "d",
		"postgresql://postgres:postgres@localhost:5432/metrics?sslmode=disable",
		"dsn for connecting to postgres")
	flag.BoolVar(&fileStorage, "fs", true, "enable fileStorage")
	flag.BoolVar(&enableHashing, "es", false, "enable hashing")
	flag.StringVar(&key, "key", "secret", "secret for hmac")

	flag.Parse()

	if v := os.Getenv("PORT"); v != "" {
		portServer = v
	}
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		storeInterval, _ = strconv.Atoi(v)
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		fileStoragePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		restore, _ = strconv.ParseBool(v)
	}
	if v := os.Getenv("DSN"); v != "" {
		dsn = v
	}
	if v := os.Getenv("FILE_STORAGE"); v != "" {
		fileStorage, _ = strconv.ParseBool(v)
	}
	if v := os.Getenv("ENABLE_HASHING"); v != "" {
		enableHashing, _ = strconv.ParseBool(v)
	}

	return portServer, storeInterval, fileStoragePath, restore, dsn, fileStorage, enableHashing, key
}

func ConfigurateAgent() (string, int, int, bool, string, int) {

	var (
		portServer     string
		pollInterval   int
		reportInterval int
		enablehashing  bool
		key            string
		rate_limit     int
	)

	flag.IntVar(&pollInterval, "pollInterval", 1, "PollIntervalValue in sec")
	flag.IntVar(&reportInterval, "reportInterval", 2, "ReportIntervalValue in sec")
	flag.StringVar(&portServer, "port", "8080", "port for server")
	flag.BoolVar(&enablehashing, "es", false, "enable hashing, default:false")
	flag.StringVar(&key, "key", "secret", "hashing key")
	flag.IntVar(&rate_limit, "rate_limit", 3, "rateLimit")
	flag.Parse()

	strPollInterval := os.Getenv("pollInterval")
	intPollinterval, err := strconv.Atoi(strPollInterval)
	if err == nil {
		pollInterval = intPollinterval
	}

	strReportInterval := os.Getenv("pollInterval")
	intReportInterval, err := strconv.Atoi(strReportInterval)

	if err == nil {
		reportInterval = intReportInterval
	}

	if envAddres := os.Getenv("PORT"); envAddres != "" {
		portServer = envAddres
	}

	if v := os.Getenv("RATE_LIMIT"); v != "" {
		rate_limit, _ = strconv.Atoi(v)
	}

	log.Printf("pollInterval: %v", pollInterval)
	log.Printf("reportInterval: %v", reportInterval)
	log.Printf("portServer: %v", portServer)
	return portServer, pollInterval, reportInterval, enablehashing, key, rate_limit
}
