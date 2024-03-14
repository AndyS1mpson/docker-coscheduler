package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

type externalClient interface {
	imageAPI
	containerAPI
	Ping(ctx context.Context) (types.Ping, error)
}

// imageAPI API для взаимодействия с docker образами
type imageAPI interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImageList(ctx context.Context, options types.ImageListOptions) ([]image.Summary, error)
}

// containerAPI API для взаимодейстия с docker контейнерами
type containerAPI interface {
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerCreate(
		ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig,
		platform *ocispec.Platform,
		containerName string,
	) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerPause(ctx context.Context, containerID string) error
	ContainerUnpause(ctx context.Context, containerID string) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerUpdate(ctx context.Context, containerID string, updateConfig container.UpdateConfig) (container.ContainerUpdateOKBody, error)
}
