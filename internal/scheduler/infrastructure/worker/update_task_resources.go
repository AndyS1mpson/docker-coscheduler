package worker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// UpdateTaskResources обновление ресурсов, выделенных задаче на воркере
func (c *Client) UpdateTaskResources(ctx context.Context, containerID string, cpuSet models.CPUSet) error {
	_, err := c.externalClient.UpdateTaskResources(ctx, &task.UpdateTaskResourcesRequest{
		ContainerId: containerID,
		CpusOpt: &task.CPUsOpt{
			From:  cpuSet.From,
			Count: cpuSet.Count,
		},
	})
	if err != nil {
		return fmt.Errorf("update task resources on worker: %w", err)
	}

	return nil
}
