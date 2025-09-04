package pkg

import (
	"log"
	"net/http"
)

func ResponceHTTP(w http.ResponseWriter, message string, statuscode int) {
	log.Printf("Responding with status: %d %s", statuscode, http.StatusText(statuscode))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statuscode)

	w.Write([]byte(message))
}
