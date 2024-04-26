package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

var experimentStrategyTaskTable = "strategy_task"

// SaveExperimentStrategyTasksTx сохранение информации о задачах, которые выполнялись в стратегии в рамках эксперимента в транзакции
func (r *Repository) SaveExperimentStrategyTasksTx(
	ctx context.Context,
	q *sqlx.Tx,
	experimentID int64,
	tasks []models.StrategyTask,
) (int64, error) {
	sql := sq.Insert(experimentStrategyTaskTable).
		Columns("experiment_id", "task_name", "task_path").
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	for _, task := range tasks {
		sql = sql.Values(experimentID, task.Name, task.FolderName)
	}

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
