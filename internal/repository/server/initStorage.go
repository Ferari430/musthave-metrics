package repository

import (
	"log"

	"github.com/Ferari430/musthave-metrics/internal/interfaces"
	fileStorage "github.com/Ferari430/musthave-metrics/internal/repository/server/file"
	inMemoryStorage "github.com/Ferari430/musthave-metrics/internal/repository/server/inMemoryDb"
	postgresStorage "github.com/Ferari430/musthave-metrics/internal/repository/server/postgres"

	"github.com/Ferari430/musthave-metrics/pkg"
)

// При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d либо при их пустых значениях,
// вернитесь последовательно к:
// хранению метрик в файле — при наличии соответствующей переменной окружения или флага командной строки;
// хранению метрик в памяти.

// возвращать наполненный интерфейс который имплементят все виды хранилищ
// todo: добавить флаг для fileStorage где проверять включена ли возможность записывать данные в файл

func InitRepository(dsn, file_storage_path string, fileStorageFlag bool) (interfaces.Repository, interfaces.FileStorage) {
	var filestorage interfaces.FileStorage
	if fileStorageFlag {
		//init file
		producer := pkg.NewProducer(file_storage_path)
		consumer := pkg.NewConsumer(file_storage_path)

		fs := fileStorage.NewFileStorage(producer, consumer)
		err := fs.Ping(file_storage_path)
		if err == nil {
			log.Println("file initialaized")
			filestorage = fs
		}
	}

	//postgres
	if dsn != "" {
		pg, err := postgresStorage.Open(dsn)
		if err == nil {

			db := postgresStorage.NewPostgresRepository(pg)
			log.Println("postgres connection initialaized")
			return db, filestorage
		}
	}

	db := inMemoryStorage.NewInMemoryStorage()
	log.Println("inMemoryStorage initialaized")
	return db, filestorage
}
