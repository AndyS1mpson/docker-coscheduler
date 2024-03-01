package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

type externalClient interface {
	Ping(ctx context.Context) (types.Ping, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}
