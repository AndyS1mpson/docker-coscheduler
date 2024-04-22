package docker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
)

// GetContainerCPUTime получение процессорного времени контейнера
func (c *Client) GetContainerCPUTime(ctx context.Context, containerID string) (time.Duration, error) {
	stats, err := c.externalClient.ContainerStats(ctx, containerID, false)
	if err != nil {
		return 0, err
	}
	defer stats.Body.Close()

	var cpuUsage time.Duration

	decoder := json.NewDecoder(stats.Body)
	for {
		var statsData types.StatsJSON
		if err := decoder.Decode(&statsData); err != nil {
			break
		}

		cpuUsage += time.Duration(statsData.CPUStats.CPUUsage.TotalUsage)
	}

	return cpuUsage, nil
}
