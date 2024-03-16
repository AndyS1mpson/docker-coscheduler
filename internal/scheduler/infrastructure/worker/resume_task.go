package worker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// ResumeTask возобновление работы задачи после паузы
func (c *Client) ResumeTask(ctx context.Context, containerID string) error {
	_, err := c.externalClient.ResumeTask(ctx, &task.ResumeTaskRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return fmt.Errorf("resume task on worker: %w", err)
	}

	return nil
}
