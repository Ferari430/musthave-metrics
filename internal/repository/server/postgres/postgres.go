package postgresStorage

import (
	"database/sql"
	"log"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func Open(connectionString string) (*sql.DB, error) {
	// conn := "postgres://postgres:postgres@localhost:5432/metrics?sslmode=disable"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("cant init postgres pool")
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (r *PostgresRepository) Ping() error {
	connectionString := "postgres://postgres:postgres@localhost:5432/postgres"

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("cant connect to db")
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Add(metrics []*models.Metrics) error {
	log.Println("add postgres")
	return nil
}

func (r *PostgresRepository) GetAll() ([]*models.Metrics, error) {
	return nil, nil
}

func (r *PostgresRepository) Update(metric *models.Metrics) (*models.Metrics, error) {
	return nil, nil
}

func (r *PostgresRepository) Metrics() []*models.Metrics {
	return nil
}

func (r *PostgresRepository) Metric(name string) (*models.Metrics, error) {
	return nil, nil
}

func (r *PostgresRepository) AddMetric(metric *models.Metrics) error {
	return nil
}

func (r *PostgresRepository) UpdateCounter(name string, value int64) (*models.Metrics, error) {
	return nil, nil
}

func (r *PostgresRepository) UpdateGauge(name string, value float64) (*models.Metrics, error) {
	return nil, nil
}

func (r *PostgresRepository) CheckExistence(name string) bool {
	return true
}
