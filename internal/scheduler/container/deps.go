package container

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/google/uuid"
)

// Strategy интерфейс, описывающий стратегии планирования
type Strategy interface {
	Execute(ctx context.Context, experimentID uuid.UUID, tasks []models.StrategyTask) (time.Duration, error)
}
