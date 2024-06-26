package models

import "time"

// ExperimentResult структура описания эксперимента запуска стратегии
type ExperimentResult struct {
	ID             int64
	IdempotencyKey string
	StrategyID     int64
	ExecutionTime  time.Duration
}
