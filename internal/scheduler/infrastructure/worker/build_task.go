package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Build сборка образа таски
func (c *Client) Build(ctx context.Context, archive models.ImageArchive, taskTitle string) (string, error) {
	res, err := c.externalClient.BuildTask(ctx, &task.BuildTaskRequest{
		ImageArchive: archive.File,
		TaskTitle:    taskTitle,
	})
	if err != nil {
		return "", err
	}

	return res.ImageId, nil
}
