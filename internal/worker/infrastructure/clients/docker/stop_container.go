package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

// StopContainer остановка docker контейнера
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	return c.externalClient.ContainerStop(ctx, containerID, container.StopOptions{})
}
