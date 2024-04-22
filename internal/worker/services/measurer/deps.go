package measurer

import (
	"context"
	"time"
)

type dockerClient interface {
	GetContainerCPUTime(ctx context.Context, containerID string) (time.Duration, error)
}
