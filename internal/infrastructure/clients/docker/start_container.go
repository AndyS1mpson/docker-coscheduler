package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

// StartContainer запускает docker контейнер
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	return c.externalClient.ContainerStart(ctx, containerID, container.StartOptions{})
}
