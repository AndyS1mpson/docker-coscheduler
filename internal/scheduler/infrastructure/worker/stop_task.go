package worker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// StopTask остановка задачи на воркере
func (c *Client) StopTask(ctx context.Context, containerID string) error {
	_, err := c.externalClient.StopTask(ctx, &task.StopTaskRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return fmt.Errorf("stop task on worker: %w", err)
	}

	return nil
}
