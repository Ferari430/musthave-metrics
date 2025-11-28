package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/internal/service/agentService"
	"github.com/Ferari430/musthave-metrics/pkg/hash"
)

type AgentSender struct {
	client      *http.Client
	service     *agentService.AgentService
	serverPort  string
	hashingFlag bool
	key         string
}

func NewAgentSender(service *agentService.AgentService, client *http.Client, serverPort string,
	hashingFlag bool, key string) *AgentSender {
	return &AgentSender{client: client, service: service,
		serverPort:  serverPort,
		hashingFlag: hashingFlag,
		key:         key,
	}
}

func (a *AgentSender) Consumer(wg *sync.WaitGroup) {
	// Эта горутина работает в течение всей жизни приложения.
	// wg.Done() не вызывается, что позволяет wg.Wait() в main блокировать программу бессрочно.
	channel := a.service.MetricsChannel()
	for metric := range channel {
		log.Printf("consumer: sending http request to server %v", metric)
		// a.SendJSONBatch(metric)
		// a.SendJsonCompressed(metric)
		a.SendJsonCompressedBatch(metric) // Compressed JSON Batch
	}
}

func (a *AgentSender) SendHTTP(metrics []*models.Metrics) {
	for _, val := range metrics {
		var url string
		switch val.MType {
		case "gauge":
			url = fmt.Sprintf("http://localhost:%v/update/%v/%v/%v", a.serverPort, "gauge", val.ID, *val.Value)
			log.Println("sending request to " + url)

		case "counter":
			url = fmt.Sprintf("http://localhost:%v/update/%v/%v/%v", a.serverPort, "counter", val.ID, *val.Delta)
			log.Println("sending request to " + url)

		}

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		// req.Header.Set("Content-Encoding", "gzip")
		resp, err := a.client.Do(req)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		resp.Body.Close()
		log.Println("response Status:", resp.Status)
	}
}

func (a *AgentSender) SendBatch(metrics []*models.Metrics) {
	log.Println("SENDING JSON")
	for _, val := range metrics {
		jsonValue, err := json.Marshal(val)
		if err != nil {
			log.Println("cant marshall metric to json")
			return
		}

		log.Printf("json: %v", string(jsonValue))
		url := "http://localhost:8080/valueJ"
		reader := bytes.NewReader(jsonValue)
		req, err := http.NewRequest("POST", url, reader)
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		// req.Header.Set("Content-Encoding", "gzip")
		resp, err := a.client.Do(req)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		resp.Body.Close()
		log.Println("response Status:", resp.Status)

	}
}

type Simple struct {
	message string
}

func (a *AgentSender) Ping() {
	log.Println("ping")
	url := "http://localhost:8080/valueJ"

	stringg := Simple{
		message: "ping",
	}

	jsonValue, err := json.Marshal(stringg)

	reader := bytes.NewReader(jsonValue)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	resp.Body.Close()
	log.Println("response Status:", resp.Status)

}

func (a *AgentSender) SendJsonCompressed(metrics []*models.Metrics) {
	log.Println("SENDING JSON")
	for _, val := range metrics {
		jsonValue, err := json.Marshal(val)
		if err != nil {
			log.Println("cant marshall metric to json")
			return
		}
		b := bytes.Buffer{}
		gzipWriter := gzip.NewWriter(&b)
		gzipWriter.Write(jsonValue)
		gzipWriter.Close() // обязателньо закрыть после записи чтобы архив отправился корректно

		log.Printf("json: %v", string(jsonValue))
		url := "http://localhost:8080/valueJ"
		reader := bytes.NewReader(b.Bytes())
		req, err := http.NewRequest("POST", url, reader)
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		req.Header.Set("Content-Encoding", "gzip")
		resp, err := a.client.Do(req)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		resp.Body.Close()
		log.Println("response Status:", resp.Status)

	}
}

func (a *AgentSender) SendJSONBatch(metrics []*models.Metrics) {
	log.Println("SENDING BATCH JSON")
	jsonValue, err := json.Marshal(metrics)
	if err != nil {
		log.Println("cant marshall metric to json")
		return
	}

	log.Printf("json: %v", string(jsonValue))
	url := "http://localhost:8080/valueJ"
	reader := bytes.NewReader(jsonValue)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
	}

	req.Header.Set("Content-Type", "application/json")

	// req.Header.Set("Content-Encoding", "gzip")
	resp, err := a.client.Do(req)
	if err != nil {
		log.Println("ERROR:", err)
	}

	resp.Body.Close()
	log.Println("response Status:", resp.Status)

}

func (a *AgentSender) SendJsonCompressedBatch(metrics []*models.Metrics) {
	log.Println("SENDING JSON")
	jsonValue, err := json.Marshal(metrics)
	if err != nil {
		log.Println("cant marshall metric to json")
		return
	}
	b := bytes.Buffer{}
	gzipWriter := gzip.NewWriter(&b)
	gzipWriter.Write(jsonValue)
	gzipWriter.Close() // обязателньо закрыть после записи чтобы архив отправился корректно

	log.Printf("json: %v", string(jsonValue))
	url := "http://localhost:8080/valueJ"
	reader := bytes.NewReader(b.Bytes())
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Content-Encoding", "gzip")

	if a.hashingFlag {
		h := hash.ComputeHash(jsonValue, a.key)

		req.Header.Set("HashSHA256", h)
		log.Println("signaturing req, hash = ", h)
	}

	resp, err := a.client.Do(req) // <-- блокируемся: ждем ответа
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	resp.Body.Close()
	log.Println("response Status:", resp.Status)

	log.Println("compressed size:", len(b.Bytes()), "original size:", len(jsonValue))

}
