package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	repository "github.com/Ferari430/musthave-metrics/internal/repository/server"

	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// chi
func main() {
	port, store_interval, file_storage_path, restore, connectionString, file_storage, hashing, key := pkg.ConfigurateServer()

	log.Printf("Server started on port %v", port)
	log.Printf("store_interval= %v, file_storage_path= %v, restore= %v", store_interval, file_storage_path, restore)

	err := app(port, file_storage_path, store_interval, restore, connectionString, file_storage, hashing, key)
	if err != nil {
		log.Fatalf("Cant run server on port %v. Error:%v", port, err)
		return
	}

}

func app(port, fPath string, store_interval int, restore bool, connectionString string, filestorage, hashingFlag bool, key string) error {
	router := chi.NewRouter()
	//logger
	logger, err := logger.InitLogger("debug")
	if err != nil {
		return err
	}

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
