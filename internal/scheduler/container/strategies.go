package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy/fcs"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

// GetRoundRobinStrategy последовательная стратегия планирования задач на узлах
func (c *Container) GetRoundRobinStrategy(nodes map[models.Node]*worker.Client) *strategy.RoundRobinStrategy[*worker.Client] {
	return container.MustOrGetNew(c.Container, func() *strategy.RoundRobinStrategy[*worker.Client] {
		return strategy.NewRoundRobinStrategy[*worker.Client](nodes, c.GetTaskHub(), c.configs.TaskInfoDelay)
	})
}

// GetFCSStrategy FCS стратегия планирования задач на узлах
func (c *Container) GetFCSStrategy(nodes map[models.Node]*worker.Client) *fcs.FCSStrategy[*worker.Client] {
	return container.MustOrGetNew(c.Container, func() *fcs.FCSStrategy[*worker.Client] {
		return fcs.NewFCSStrategy[*worker.Client](
			nodes,
			c.GetTaskHub(),
			c.configs.TaskInfoDelay,
			c.configs.TaskCombinationNum,
			c.configs.MeasurementTime,
		)
	})
}
