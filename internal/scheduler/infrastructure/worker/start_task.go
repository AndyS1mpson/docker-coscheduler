package worker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// StartTask запуск задачи на воркере
func (c *Client) StartTask(ctx context.Context, containerID string) error {
	_, err := c.externalClient.StartTask(ctx, &task.StartTaskRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return fmt.Errorf("start task on worker: %w", err)
	}

	return nil
}
