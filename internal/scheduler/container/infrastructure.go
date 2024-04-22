package container

import (
	imageHub "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/task_hub"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

func (c *Container) GetTaskHub() *imageHub.Hub {
	return container.MustOrGetNew(c.Container, func() *imageHub.Hub {
		return imageHub.NewHub(
			c.configs.TaskDir,
		)
	})
}
