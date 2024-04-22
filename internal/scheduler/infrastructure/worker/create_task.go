package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

var errContainerNotCreated = errors.New("task container not created")

func (c *Client) CreateTask(ctx context.Context, imageID string, cpuSet *models.CPUSet) (string, error) {
	var taskResources *task.CPUsOpt

	if cpuSet != nil {
		taskResources = &task.CPUsOpt{
			From:  cpuSet.From,
			Count: cpuSet.Count,
		}
	}

	res, err := c.externalClient.CreateTask(ctx, &task.CreateTaskRequest{
		ImageId: imageID,
		CpusOpt: taskResources,
	})
	if err != nil {
		return "", fmt.Errorf("create task on worker: %w", err)
	}

	if res.Status != string(models.ContainerStateCreated) {
		return "", errContainerNotCreated
	}

	return res.ContainerId, nil
}
