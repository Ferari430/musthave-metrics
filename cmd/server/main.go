package main

import (
	"fmt"
	"log"

	"github.com/Ferari430/musthave-metrics/cmd/server/app"
	"github.com/Ferari430/musthave-metrics/pkg"
)

// chi
func main() {
	port, store_interval, file_storage_path, restore, connectionString, file_storage, hashing, key := pkg.ConfigurateServer()

	fmt.Print("starting server....\n")
	err := app.StartApp(port, file_storage_path, store_interval, restore, connectionString, file_storage, hashing, key)
	if err != nil {
		log.Fatalf("Cant run server on port %v. Error:%v", port, err)
		return
	}

}
