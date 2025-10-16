package fileStorage

import (
	"errors"
	"log"
	"os"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg"
)

type FileStorage struct {
	Producer *pkg.Producer
	Consumer *pkg.Consumer
}

func NewFileStorage(producer *pkg.Producer, consumer *pkg.Consumer) *FileStorage {
	return &FileStorage{

		Producer: producer,
		Consumer: consumer,
	}
}

func (f *FileStorage) GetAll() ([]*models.Metrics, error) {
	_, err := f.Consumer.File.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var metrics []*models.Metrics
	err = f.Consumer.Decoder.Decode(&metrics)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (f *FileStorage) Ping(fpath string) error {

	file, err := os.OpenFile(fpath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func (f *FileStorage) Add(metrics []*models.Metrics) error {
	err := f.Producer.ClearFile()
	if err != nil {
		return err
	}
	err = f.Producer.Encoder.Encode(metrics)
	if err != nil {
		log.Println("Cant write metrics on file")
		return err
	}

	log.Println("Metrics written on file")
	return nil
}

func (f *FileStorage) Restore() ([]*models.Metrics, error) {

	metrics, err := f.Consumer.Restore()
	if err != nil {
		return nil, errors.New("cant restore data from file")
	}

	return metrics, nil
}

////////////////////////////////////////////////////////////////////
