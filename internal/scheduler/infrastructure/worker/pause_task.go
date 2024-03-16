package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// PauseTask остановка контейнера с задачей на воркере
func (c *Client) PauseTask(ctx context.Context, containerID string) error {
	_, err := c.externalClient.PauseTask(ctx, &task.PauseTaskRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return err
	}

	return nil
}
