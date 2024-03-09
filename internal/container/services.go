package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/services/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

// GetTaskService сервис, содержащий бизнес-логику работы с задачами
func (c *Container) GetTaskService() *task.Service {
	return container.MustOrGetNew(c.Container, func() *task.Service {
		return task.NewService(
			c.getDockerClient(),
			c.getImageHub(),
			*c.configs.NodeConfig,
		)
	})
}
