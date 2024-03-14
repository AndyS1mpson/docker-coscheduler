package docker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/docker/docker/api/types/container"
)

// UpdateContainer обновление конфигурации контейнера
func (c *Client) UpdateContainer(ctx context.Context, containerID string, cpuSet models.CPUSet) error {
	_, err := c.externalClient.ContainerUpdate(ctx, containerID, container.UpdateConfig{
		Resources: container.Resources{
			CpusetCpus: cpuSet.AsString(),
		},
	})

	return err
}
