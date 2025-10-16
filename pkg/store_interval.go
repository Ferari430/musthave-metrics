package pkg

import (
	"encoding/json"
	"log"
	"os"

	models "github.com/Ferari430/musthave-metrics/internal/model"
)

type Producer struct {
	File    *os.File
	Encoder *json.Encoder
}

func NewProducer(fPath string) *Producer {
	file, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}

	return &Producer{File: file,
		Encoder: json.NewEncoder(file),
	}
}

func (p *Producer) Close() error {
	return p.File.Close()
}

func (p *Producer) WriteMetric(metrics []*models.Metrics) error {
	err := p.ClearFile()
	if err != nil {
		return err
	}

	err = p.Encoder.Encode(metrics)
	if err != nil {
		log.Println("Cant write metrics on file")
		return err
	}

	log.Println("Metrics written on file")
	return nil
}

func (p *Producer) ClearFile() error {

	err := p.File.Truncate(0)
	if err != nil {
		return err
	}

	_, err = p.File.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}

type Consumer struct {
	File    *os.File
	Decoder *json.Decoder
}

func (c *Consumer) Close() error {
	return c.File.Close()
}

func NewConsumer(fPath string) *Consumer {
	file, err := os.OpenFile(fPath, os.O_RDONLY, 0666)
	if err != nil {
		return nil
	}

	return &Consumer{
		File:    file,
		Decoder: json.NewDecoder(file),
	}
}

func (c *Consumer) Restore() ([]*models.Metrics, error) {
	var metrics []*models.Metrics
	err := c.Decoder.Decode(&metrics)
	if err != nil {

		return nil, err
	}

	log.Println("Metrics restored from file")

	return metrics, nil
}

func (c *Consumer) GetAll() ([]*models.Metrics, error) {
	var metrics []*models.Metrics

	err := c.Decoder.Decode(&metrics)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
