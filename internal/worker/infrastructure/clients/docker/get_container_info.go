package docker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// GetContainerInfo получает информацию о состоянии контейнера
func (c *Client) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	result, err := c.externalClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("get container info: %w", err)
	}

	return &models.ContainerInfo{
		ID:       containerID,
		State:    models.ContainerState(result.State.Status),
		ExitCode: int64(result.State.ExitCode),
	}, nil
}
