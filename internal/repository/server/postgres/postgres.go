package postgresStorage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	models "github.com/Ferari430/musthave-metrics/internal/model"
	"github.com/Ferari430/musthave-metrics/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger *logger.Logger
	DB     *sql.DB
}

func NewPostgresRepository(db *sql.DB, logger *logger.Logger) *PostgresRepository {
	return &PostgresRepository{DB: db,
		logger: logger}
}

func Open(connectionString string, logger *logger.Logger) (*sql.DB, error) {
	op := "postgres.Open"
	// conn := "postgres://postgres:postgres@localhost:5432/metrics?sslmode=disable"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		logger.Debug("cant open postgres connection pool", zap.String("operation", op), zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (r *PostgresRepository) Ping() error {
	op := "Postrges.Ping"

	err := r.DB.Ping()

	if err != nil {
		r.logger.Debug("cant ping db", zap.String("operation", op), zap.Error(err))
	}

	r.logger.Debug("connection successes", zap.String("operation", op), zap.Error(nil))

	return nil
}

func (r *PostgresRepository) Add(metrics []*models.Metrics) error {
	op := "Postgres.Add"

	if len(metrics) == 0 {
		r.logger.Debug("recieved metrics with len=0", zap.String("operation", op), zap.Error(errors.New("zero len")))
		return nil
	}

	placeholders := make([]string, 0, len(metrics))
	args := make([]interface{}, 0, len(metrics)*5)

	for i, m := range metrics {
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		args = append(args, m.ID, m.MType, m.Delta, m.Value, m.Hash)
	}

	query := fmt.Sprintf(`
        INSERT INTO metrics (id, mtype, delta, value, hash) 
        VALUES %s
        ON CONFLICT (id) DO UPDATE SET
            delta = CASE 
                WHEN excluded.mtype = 'Counter' THEN metrics.delta + excluded.delta
                ELSE excluded.delta
            END,
            value = excluded.value,
            hash = excluded.hash
        WHERE metrics.mtype = excluded.mtype`,
		strings.Join(placeholders, ", "))

	result, err := r.DB.Exec(query, args...)
	if err != nil {

		r.logger.Debug("error during insert", zap.String("operation", op), zap.Error(err))

		return err
	}

	rowsAffected, _ := result.RowsAffected()

	r.logger.Debug("Insert successful", zap.String("operation", op), zap.Int64("rows_affected", rowsAffected))

	return err
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
