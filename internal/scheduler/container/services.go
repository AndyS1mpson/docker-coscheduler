package container

import (
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/repository"
	nodeResourcesCache "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/node/cache/resources"
	nodeTasksSpeedCache "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/node/cache/task_speed"
	strategiesCache "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/strategy/cache"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

func (c *Container) getRepository() *repository.Repository {
	return container.MustOrGetNew(c.Container, func() *repository.Repository {
		return repository.New(
			c.getPgConnection(),
		)
	})
}

func (c *Container) getNodeSpeedCache(nodes []models.Node) *nodeTasksSpeedCache.Cache {
	return container.MustOrGetNew(c.Container, func() *nodeTasksSpeedCache.Cache {
		return nodeTasksSpeedCache.NewCache(nodes, c.configs.FCNTaskNum)
	})
}

func (c *Container) getStrategiesCache() *strategiesCache.Cache {
	return container.MustOrGetNew(c.Container, func() *strategiesCache.Cache {
		cache := strategiesCache.NewCache(c.getRepository(), c.configs.StrategyLoadingBatchSize, c.configs.CacheLoadingReloadInterval)

		if err := cache.Load(c.Ctx()); err != nil {
			panic(fmt.Errorf("load strategies: %w", err))
		}

		cache.StartLoading(c.Ctx())

		return cache
	})
}

func (c *Container) getNodeResourcesCache(nodes map[models.Node]*worker.Client) *nodeResourcesCache.Cache[*worker.Client] {
	return container.MustOrGetNew(c.Container, func() *nodeResourcesCache.Cache[*worker.Client] {
		cache := nodeResourcesCache.NewCache(nodes, c.configs.NodeResourcesRequestReloadInterval)

		cache.StartLoading(c.Ctx())

		return cache
	})
}
