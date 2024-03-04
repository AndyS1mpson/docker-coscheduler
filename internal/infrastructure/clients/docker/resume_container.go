package docker

import "context"

// UnpauseContainer продолжает выполнение контейнера
func (c *Client) UnpauseContainer(ctx context.Context, containerID string) error {
	return c.externalClient.ContainerUnpause(ctx, containerID)
}
