package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy/fcn"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy/fcs"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
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

// GetFCNStrategy FCN стратегия планирования задач на узлах
func (c *Container) GetFCNStrategy(nodes map[models.Node]*worker.Client) *fcn.FCNStrategy[*worker.Client] {
	return container.MustOrGetNew(c.Container, func() *fcn.FCNStrategy[*worker.Client] {
		return fcn.NewFCNStrategy[*worker.Client](
			nodes,
			c.getNodeSpeedCache(slices.Keys(nodes)),
			c.GetTaskHub(),
			c.configs.FCNTaskNum,
			c.configs.TaskInfoDelay,
			c.configs.MeasurementTime,
		)
	})
}
