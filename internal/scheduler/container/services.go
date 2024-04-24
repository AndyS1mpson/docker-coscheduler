package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/cache"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

func (c *Container) getNodeSpeedCache(nodes []models.Node) *cache.Cache {
	return container.MustOrGetNew(c.Container, func() *cache.Cache {
		return cache.NewCache(nodes, c.configs.FCNTaskNum)
	})
}
