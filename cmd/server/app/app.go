package app

import (
	"net/http"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	repository "github.com/Ferari430/musthave-metrics/internal/repository/server"
	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func StartApp(port, fPath string, store_interval int, restore bool, connectionString string, filestorage, hashingFlag bool, key string) error {
	router := chi.NewRouter()
	//logger
	logger, err := logger.InitLogger("debug")
	if err != nil {
		return err
	}

	logger.Debug(
		"server started with follow configurations:",
		zap.String("port", port),
		zap.String("file_path", fPath),
		zap.Int("store_interval", store_interval),
		zap.Bool("restore", restore),
		zap.String("connection_string", connectionString),
		zap.Bool("file_storage", filestorage),
		zap.Bool("hashing_flag", hashingFlag),
		zap.String("key", key),
	)

	repo, file := repository.InitRepository(connectionString, fPath, filestorage, logger)

	//ticker
	ticker := time.NewTicker(time.Second * time.Duration(store_interval))

	serviceServer := service.NewServiceServer(ticker, repo, file, logger)

	cfg := handler.HandlerConfig{Hashing: hashingFlag, Key: key}
	newServerHandlerDeps := handler.ServerHandlerDeps{Service: serviceServer, Config: cfg}
	handler.NewServerHandler(router, newServerHandlerDeps, logger)

	if restore {
		serviceServer.RestoreData()
	}

	return http.ListenAndServe(port, router)
}
