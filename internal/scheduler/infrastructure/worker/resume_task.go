package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// ResumeTask возобновление работы задачи после паузы
func (c *Client) ResumeTask(ctx context.Context, containerID string) error {
	_, err := c.externalClient.ResumeTask(ctx, &task.ResumeTaskRequest{
		ContainerId: containerID,
	})
	if err != nil {
		return err
	}

	return nil
}
