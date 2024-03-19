package worker

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Build сборка образа таски
func (c *Client) BuildTask(ctx context.Context, archive models.ImageArchive, taskTitle string) (string, error) {
	res, err := c.externalClient.BuildTask(ctx, &task.BuildTaskRequest{
		ImageArchive: archive.File,
		TaskTitle:    taskTitle,
	})
	if err != nil {
		return "", fmt.Errorf("build task on worker: %w", err)
	}

	return res.ImageId, nil
}
