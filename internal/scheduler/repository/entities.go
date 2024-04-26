package repository

import "time"

type ExperimentResult struct {
	ID             int64         `db:"id"`
	IdempotencyKey string        `db:"idempotency_key"`
	StrategyName   string        `db:"strategy_name"`
	ExecutionTime  time.Duration `db:"execution_time"`
}
