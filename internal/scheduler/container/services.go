package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/repository"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/cache"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

func (c *Container) getRepository() *repository.Repository {
	return container.MustOrGetNew(c.Container, func() *repository.Repository {
		return repository.New(
			c.getPgConnection(),
		)
	})
}

func (c *Container) getNodeSpeedCache(nodes []models.Node) *cache.Cache {
	return container.MustOrGetNew(c.Container, func() *cache.Cache {
		return cache.NewCache(nodes, c.configs.FCNTaskNum)
	})
}
