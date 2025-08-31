package main

import (
	"log"
	"net/http"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	"github.com/Ferari430/musthave-metrics/internal/repository"
	"github.com/Ferari430/musthave-metrics/internal/service"
)

func main() {

	port := ":8080"
	log.Printf("Server started on port %v", port)
	err := app(port)
	if err != nil {
		log.Fatalf("Cant run server on port %v", port)
		return
	}

}

func app(port string) error {
	router := http.NewServeMux()

	//repo
	InMemoryDB := repository.NewInMemoryRepo()

	service := service.NewServiceServer(InMemoryDB)
	newServerHandlerDeps := handler.ServerHandlerDeps{Service: service}
	handler.NewServerHandler(router, newServerHandlerDeps)

	return http.ListenAndServe(port, router)
}
