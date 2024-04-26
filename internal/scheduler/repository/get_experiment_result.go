package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/experiment"
)

// GetExperimentResultsTx получение результатов экспериментов по фильтрам в транзакции
func (r *Repository) GetExperimentResults(ctx context.Context, q *sqlx.Tx, criteria experiment.SearchCriteria) ([]models.ExperimentResult, error) {
	sql := sq.Select("*").From(tableName).PlaceholderFormat(sq.Dollar)

	if criteria.IDs != nil {
		sql = sql.Where(sq.Eq{"id": criteria.IDs})
	}

	if criteria.IdempotencyKeys != nil {
		sql = sql.Where(sq.Eq{"idempotency_key": criteria.IdempotencyKeys})
	}

	if criteria.StrategyNames != nil {
		sql = sql.Where(sq.Eq{"strategy_name": criteria.StrategyNames})
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return nil, err
	}

	var experimentsResults []ExperimentResult

	err = q.SelectContext(ctx, &experimentsResults, raw, args...)
	if err != nil {
		return nil, fmt.Errorf("find experiments results failed: %w", err)
	}

	experiments := make([]models.ExperimentResult, 0, len(experimentsResults))

	for _, entity := range experimentsResults {
		experiments = append(experiments, models.ExperimentResult{
			ID:             entity.ID,
			IdempotencyKey: entity.IdempotencyKey,
			StrategyName:   models.StrategyName(entity.StrategyName),
			ExecutionTime:  entity.ExecutionTime,
		})
	}

	return experiments, nil
}
