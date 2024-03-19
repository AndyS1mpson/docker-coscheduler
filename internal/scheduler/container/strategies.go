package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

func (c *Container) GetSequentialStrategy(nodes []*worker.Client) *strategy.TimeDecorator {
	return container.MustOrGetNew(c.Container, func() *strategy.TimeDecorator {
		return strategy.NewTimeDecorator(
			strategy.NewSequentialStrategy[*worker.Client](nodes, c.GetTaskHub(), c.configs.TaskInfoDelay),
		)
	})
}
