package docker

import "context"

// PauseContainer остановка выполнения контейнера
func (c *Client) PauseContainer(ctx context.Context, containerID string) error {
	return c.externalClient.ContainerPause(ctx, containerID)
}
