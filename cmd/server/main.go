package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	repository "github.com/Ferari430/musthave-metrics/internal/repository/server"
	fileStorage "github.com/Ferari430/musthave-metrics/internal/repository/server/file"

	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// chi
func main() {
	port, store_interval, file_storage_path, restore, connectionString := pkg.ConfigurateServer()

	log.Printf("Server started on port %v", port)
	log.Printf("store_interval= %v, file_storage_path= %v, restore= %v", store_interval, file_storage_path, restore)

	err := app(port, file_storage_path, store_interval, restore, connectionString)
	if err != nil {
		log.Fatalf("Cant run server on port %v. Error:%v", port, err)
		return
	}

}
func app(port, fPath string, store_interval int, restore bool, connectionString string) error {
	router := chi.NewRouter()
	//logger
	err := logger.InitLogger("debug")
	if err != nil {
		return err
	}
	producer := pkg.NewProducer(fPath)
	consumer := pkg.NewConsumer(fPath)
	f := fileStorage.NewFileStorage(producer, consumer)

	repo := repository.InitRepository(connectionString, fPath)
	err = repo.Ping(connectionString)
	if err != nil {
		log.Println(err)
	}

	//ticker
	ticker := time.NewTicker(time.Second * time.Duration(store_interval))

	serviceServer := service.NewServiceServer(ticker, repo, f)
	newServerHandlerDeps := handler.ServerHandlerDeps{Service: serviceServer}
	handler.NewServerHandler(router, newServerHandlerDeps)

	//restore
	// if restore {
	// 	serviceServer.RestoreData()
	// }

	return http.ListenAndServe(port, router)
}
