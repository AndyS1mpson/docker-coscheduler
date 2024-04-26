package models

import "time"

// ExperimentResult структура описания эксперимента запуска стратегии
type ExperimentResult struct {
	ID             int64
	IdempotencyKey string
	StrategyName   StrategyName
	ExecutionTime  time.Duration
}
