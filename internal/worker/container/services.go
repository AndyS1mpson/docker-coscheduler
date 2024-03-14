package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/worker/services/task"
)

// GetTaskService сервис, содержащий бизнес-логику работы с задачами
func (c *Container) GetTaskService() *task.Service {
	return container.MustOrGetNew(c.Container, func() *task.Service {
		return task.NewService(
			c.getDockerClient(),
			c.GetImageHub(),
			c.configs.Node.Host,
			c.configs.Node.Port,
			c.configs.Node.CPUNums,
		)
	})
}
