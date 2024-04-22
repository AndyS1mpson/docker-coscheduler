package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"google.golang.org/protobuf/types/known/durationpb"
)

// MeasureTaskSpeed измеренение времени выполнения задачи
func (c *Client) MeasureTaskSpeed(
	ctx context.Context,
	containerID string,
	cpuSet models.CPUSet,
	duration time.Duration,
) (time.Duration, error) {
	res, err := c.externalClient.MeasureTaskSpeed(ctx, &task.MeasureTaskSpeedRequest{
		Duration:    durationpb.New(duration),
		ContainerId: containerID,
		CpusOpt: &task.CPUsOpt{
			From:  cpuSet.From,
			Count: cpuSet.Count,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("measure task speed on worker: %w", err)
	}

	return res.Time.AsDuration(), nil
}
