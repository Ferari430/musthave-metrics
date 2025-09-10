package main

import (
	"log"
	"net/http"

	"github.com/Ferari430/musthave-metrics/internal/handler"
	"github.com/Ferari430/musthave-metrics/internal/repository"
	"github.com/Ferari430/musthave-metrics/internal/service"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/go-chi/chi/v5"
)

// chi
func main() {
	port := pkg.ConfigurateServer()
	log.Printf("Server started on port %v", port)
	err := app(port)
	if err != nil {
		log.Fatalf("Cant run server on port %v. Error:%v", port, err)
		return
	}

}

func app(port string) error {
	router := chi.NewRouter()

	//logger
	pkg.InitLogger("debug")
	//repo
	InMemoryDB := repository.NewInMemoryRepo()

	service := service.NewServiceServer(InMemoryDB)
	newServerHandlerDeps := handler.ServerHandlerDeps{Service: service}
	handler.NewServerHandler(router, newServerHandlerDeps)

	return http.ListenAndServe(port, router)
}
