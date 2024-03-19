package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
)

type TimeDecorator struct {
	strategy strategy
}

// NewTimeDecorator конструктор создания TimeDecorator
func NewTimeDecorator(strategy strategy) *TimeDecorator {
	return &TimeDecorator{
		strategy: strategy,
	}
}

// Execute
func (s *TimeDecorator) Execute(ctx context.Context, tasks []models.StrategyTask) {
	start := time.Now()

	s.strategy.Execute(ctx, tasks)

	log.Info(fmt.Sprintf("strategy execution time: %s", time.Since(start)), log.Data{})
}
