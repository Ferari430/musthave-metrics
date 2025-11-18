package fileStorage

import (
	"errors"
	"os"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	"go.uber.org/zap"
)

type FileStorage struct {
	Producer *pkg.Producer
	Consumer *pkg.Consumer
	logger   *logger.Logger
}

func NewFileStorage(producer *pkg.Producer, consumer *pkg.Consumer, log *logger.Logger) *FileStorage {
	return &FileStorage{
		Producer: producer,
		Consumer: consumer,
		logger:   log,
	}
}

func (f *FileStorage) GetAll() ([]*models.Metrics, error) {
	op := "FileStorage.GetAll"

	_, err := f.Consumer.File.Seek(0, 0)
	if err != nil {
		f.logger.Error("failed to seek file",
			zap.String("operation", op),
			zap.Error(err),
		)
		return nil, err
	}

	var metrics []*models.Metrics
	err = f.Consumer.Decoder.Decode(&metrics)
	if err != nil {
		f.logger.Error("failed to decode metrics",
			zap.String("operation", op),
			zap.Error(err),
		)
		return nil, err
	}

	f.logger.Debug("metrics decoded from file",
		zap.String("operation", op),
		zap.Int("count", len(metrics)),
	)

	return metrics, nil
}

func (f *FileStorage) Ping(fpath string) error {
	op := "FileStorage.Ping"

	file, err := os.OpenFile(fpath, os.O_RDWR, 0644)
	if err != nil {
		f.logger.Error("failed to open file",
			zap.String("operation", op),
			zap.String("path", fpath),
			zap.Error(err),
		)
		return err
	}
	defer file.Close()

	f.logger.Debug("file storage reachable",
		zap.String("operation", op),
		zap.String("path", fpath),
	)

	return nil
}

func (f *FileStorage) Add(metrics []*models.Metrics) error {
	op := "FileStorage.Add"

	err := f.Producer.ClearFile()
	if err != nil {
		f.logger.Error("failed to clear file",
			zap.String("operation", op),
			zap.Error(err),
		)
		return err
	}

	err = f.Producer.Encoder.Encode(metrics)
	if err != nil {
		f.logger.Error("failed to encode metrics into file",
			zap.String("operation", op),
			zap.Error(err),
		)
		return err
	}

	f.logger.Info("metrics written to file",
		zap.String("operation", op),
		zap.Int("count", len(metrics)),
	)

	return nil
}

func (f *FileStorage) Restore() ([]*models.Metrics, error) {
	op := "FileStorage.Restore"

	metrics, err := f.Consumer.Restore()
	if err != nil {
		f.logger.Error("failed to restore from file",
			zap.String("operation", op),
			zap.Error(err),
		)
		return nil, errors.New("cant restore data from file")
	}

	f.logger.Info("restored metrics from file",
		zap.String("operation", op),
		zap.Int("count", len(metrics)),
	)

	return metrics, nil
}
