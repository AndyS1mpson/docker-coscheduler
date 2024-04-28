package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/strategy"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/sql"
)

// Cache кэш стратегий планирования
type Cache struct {
	mu               sync.RWMutex
	strategiesByName map[models.StrategyName]models.Strategy
	strategiesByID   map[int64]models.Strategy

	repo           repository
	batchSize      int64
	reloadInterval time.Duration
}

// NewCache конструктор для Cache
func NewCache(repo repository, batchSize int64, reloadInterval time.Duration) *Cache {
	return &Cache{
		repo:             repo,
		batchSize:        batchSize,
		reloadInterval:   reloadInterval,
		strategiesByName: make(map[models.StrategyName]models.Strategy),
		strategiesByID:   make(map[int64]models.Strategy),
	}
}

// StartLoading запускает асинхронную загрузку стратегий в кэш
func (c *Cache) StartLoading(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(c.reloadInterval):
				if err := c.Load(ctx); err != nil {
					log.Error(err, log.Data{})
				}
			}
		}
	}()
}

// Load загрузка стратегий в кэш
func (c *Cache) Load(ctx context.Context) error {
	strategies := make([]models.Strategy, 0, c.batchSize)

	var lastReceivedID int64

	for {
		batch, err := c.repo.FindStrategies(ctx, getPageQuery(lastReceivedID, c.batchSize))
		if err != nil {
			return fmt.Errorf("find strategies: %w", err)
		}

		if len(batch) == 0 {
			break
		}

		lastReceivedID = batch[len(batch)-1].ID
		strategies = append(strategies, batch...)

		if int64(len(batch)) < c.batchSize {
			break
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, strategy := range strategies {
		c.strategiesByID[strategy.ID] = strategy
		c.strategiesByName[strategy.Name] = strategy
	}

	return nil
}

// GetByName получить стратегию по названию
func (c *Cache) GetByName(ctx context.Context, name models.StrategyName) *models.Strategy {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c, ok := c.strategiesByName[name]; ok {
		return &c
	}

	return nil
}

// GetByID получить стратегию по id
func (c *Cache) GetByID(ctx context.Context, id int64) *models.Strategy {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c, ok := c.strategiesByID[id]; ok {
		return &c
	}

	return nil
}

func getPageQuery(lastReceivedID int64, batchSize int64) strategy.SearchCriteria {
	return strategy.SearchCriteria{
		PageQuery: &sql.PageQuery{
			LastIDQuery: &sql.LastIDPageQuery{
				LastRecievedID: lastReceivedID,
				PageSize:       batchSize,
			},
			OrderBy: sql.OrderBy{
				Columns: []sql.OrderByColumn{
					{
						Column:    "id",
						Direction: sql.OrderDirectionAsc,
					},
				},
			},
		},
	}
}
