package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy/fcn"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy/fcs"
	roundRobin "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/strategy/round_robin"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

// GetRoundRobinStrategy последовательная стратегия планирования задач на узлах
func (c *Container) GetRoundRobinStrategy(nodes map[models.Node]*worker.Client) *roundRobin.RoundRobinStrategy[*worker.Client] {
	return container.MustOrGetNew(c.Container, func() *roundRobin.RoundRobinStrategy[*worker.Client] {
		return roundRobin.NewRoundRobinStrategy[*worker.Client](
			c.GetStorage(),
			c.getRepository(),
			c.getStrategiesCache(),
			nodes, c.GetTaskHub(),
			c.configs.TaskInfoDelay,
		)
	})
}

// GetFCSStrategy FCS стратегия планирования задач на узлах
func (c *Container) GetFCSStrategy(nodes map[models.Node]*worker.Client) *fcs.FCSStrategy[*worker.Client] {
	return container.MustOrGetNew(c.Container, func() *fcs.FCSStrategy[*worker.Client] {
		return fcs.NewFCSStrategy[*worker.Client](
			c.GetStorage(),
			c.getRepository(),
			c.getStrategiesCache(),
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
			c.GetStorage(),
			c.getRepository(),
			c.getStrategiesCache(),
			nodes,
			c.getNodeSpeedCache(slices.Keys(nodes)),
			c.GetTaskHub(),
			c.configs.FCNTaskNum,
			c.configs.TaskInfoDelay,
			c.configs.MeasurementTime,
		)
	})
}
