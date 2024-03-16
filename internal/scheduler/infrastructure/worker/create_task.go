package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

var containerNotCreatedErr = errors.New("task container not created")

func (c *Client) CreateTask(ctx context.Context, imageID string, cpuSet models.CPUSet) (string, error) {
	res, err := c.externalClient.CreateTask(ctx, &task.CreateTaskRequest{
		ImageId: imageID,
		CpusOpt: &task.CPUsOpt{
			From:  cpuSet.From,
			Count: cpuSet.Count,
		},
	})
	if err != nil {
		return "", fmt.Errorf("create task on worker: %w", err)
	}

	if res.Status != string(models.ContainerStateCreated) {
		return "", containerNotCreatedErr
	}

	return res.ContainerId, nil
}
