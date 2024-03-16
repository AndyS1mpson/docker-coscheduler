package worker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// GetContainerInfo получение информации о контейнере с воркера
func (c *Client) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	res, err := c.externalClient.GetContainerInfo(ctx, &task.GetContainerInfoRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return nil, fmt.Errorf("get container info from worker: %w", err)
	}

	return &models.ContainerInfo{
		ID:       res.Id,
		State:    models.ContainerState(res.State),
		ExitCode: res.ExitCode,
	}, nil
}
