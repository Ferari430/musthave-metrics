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

func ConfigurateAgent() (string, int64, int64, bool, string) {

	var (
		portServer     string
		pollInterval   int64
		reportInterval int64
		enablehashing  bool
		key            string
	)

	flag.Int64Var(&pollInterval, "pollInterval", 2, "PollIntervalValue in sec")
	flag.Int64Var(&reportInterval, "reportInterval", 3, "ReportIntervalValue in sec")
	flag.StringVar(&portServer, "port", "8080", "port for server")
	flag.BoolVar(&enablehashing, "es", false, "enable hashing, default:false")
	flag.StringVar(&key, "key", "secret", "hashing key")
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
	return portServer, pollInterval, reportInterval, enablehashing, key
}
