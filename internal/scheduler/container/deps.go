package container

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Strategy интерфейс, описывающий стратегии планирования
type Strategy interface {
	Execute(ctx context.Context, tasks []models.StrategyTask) (time.Duration, error)
}
