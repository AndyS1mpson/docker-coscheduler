package repository

import "time"

// Strategy описание таблицы strategy в базе данных
type Strategy struct {
	ID           int64  `db:"id"`
	StrategyName string `db:"strategy_name"`
}

// ExperimentResult описание таблицы experiment в базе данных
type ExperimentResult struct {
	ID             int64         `db:"id"`
	IdempotencyKey string        `db:"idempotency_key"`
	StrategyID     int64         `db:"strategy_id"`
	ExecutionTime  time.Duration `db:"execution_time"`
}
