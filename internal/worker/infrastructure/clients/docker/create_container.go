package docker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/docker/docker/api/types/container"
)

// CreateContainer создание docker контейнера на основе docker образа
func (c *Client) CreateContainer(ctx context.Context, imageID string, cpuSet *models.CPUSet, containerName string) (string, error) {
	config := &container.Config{
		Image: imageID,
	}

	var resources container.Resources

	if cpuSet != nil {
		resources.CpusetCpus = cpuSet.AsString()
	}

	containerConfig := &container.HostConfig{
		Resources: resources,
	}

	res, err := c.externalClient.ContainerCreate(ctx, config, containerConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
