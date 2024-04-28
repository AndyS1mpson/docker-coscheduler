package cache

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/strategy"
)

type repository interface {
	FindStrategies(ctx context.Context, criteria strategy.SearchCriteria) ([]models.Strategy, error)
}
