package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// WaitForTask ожидание выполнения задачи на ноде
func (c *Client) WaitForTask(ctx context.Context, taskID string, delay time.Duration) error {
	isRunning := true

	for isRunning {
		info, err := c.GetTaskInfo(ctx, taskID)
		if err != nil {
			return fmt.Errorf("get task info: %w", err)
		}

		if info.State == models.ContainerStateRunning {
			time.Sleep(delay)
			continue
		} else if info.State == models.ContainerStateExited && info.ExitCode == 0 {
			isRunning = false
		} else {
			return fmt.Errorf("task crashed with status: %w and exit code: %d", info.State, info.ExitCode)
		}
	}

	return nil
}
