package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

var tableName = "experiment"

// SaveExperimentResultTx сохранение результатов экспериментов в транзакции
func (r *Repository) SaveExperimentResultTx(ctx context.Context, q *sqlx.Tx, result models.ExperimentResult) (int64, error) {
	sql := sq.Insert(tableName).
		Columns("idempotency_key", "strategy_name", "execution_time").
		Values(result.IdempotencyKey, result.StrategyName, result.ExecutionTime).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	raw, args, err := sql.ToSql()
	if err != nil {

		return 0, err
	}

	var id int64

	err = q.GetContext(ctx, &id, raw, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to exec %s insert: %w", tableName, err)
	}

	return id, nil
}
