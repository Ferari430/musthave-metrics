package postgres

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
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("cant connect to db")
	}
	log.Println("Postgres connected")
	return db, nil
}

func (r *PostgresRepository) Ping(connectionString string) error {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("cant connect to db")
	}

	return db.Ping()
}

func (r *PostgresRepository) Add(metrics *models.Metrics) error {
	return nil
}
