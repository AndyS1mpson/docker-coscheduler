package repository

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/strategy"

	sq "github.com/Masterminds/squirrel"
)

var strategyTableName = "strategy"

// FindStrategies поиск стратегий по фильтрам
func (r *Repository) FindStrategies(ctx context.Context, criteria strategy.SearchCriteria) ([]models.Strategy, error) {
	sql := sq.Select("*").From(strategyTableName).PlaceholderFormat(sq.Dollar)

	if criteria.IDs != nil {
		sql = sql.Where(sq.Eq{"id": criteria.IDs})
	}

	if criteria.Names != nil {
		sql = sql.Where(sq.Eq{"name": criteria.Names})
	}

	raw, args, err := sql.ToSql()
	if err != nil {
		return nil, err
	}

	var entities []Strategy

	err = r.db.SelectContext(ctx, &entities, raw, args...)
	if err != nil {
		return nil, fmt.Errorf("find strategies: %w", err)
	}

	strategies := make([]models.Strategy, len(entities))

	for _, entity := range entities {
		strategies = append(strategies, models.Strategy{
			ID:   entity.ID,
			Name: models.StrategyName(entity.StrategyName),
		})
	}

	return strategies, nil
}
