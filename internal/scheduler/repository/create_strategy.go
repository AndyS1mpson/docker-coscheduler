package repository

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"

	sq "github.com/Masterminds/squirrel"
)

// CreateStrategy создание стратегии в базе данных
func (r *Repository) CreateStrategy(ctx context.Context, strategy models.Strategy) (int64, error) {
	sql := sq.Insert(strategyTableName).
		Columns("name").
		Values(strategy.Name).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	raw, args, err := sql.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64

	err = r.db.GetContext(ctx, &id, raw, args...)
	if err != nil {
		return 0, fmt.Errorf("create promocode execution failed: %w", err)
	}

	return id, nil
}
