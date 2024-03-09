package container

import (
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/infrastructure/clients/docker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

func (c *Container) getDockerClient() *docker.Client {
	return container.MustOrGetNew(c.Container, func() *docker.Client {
		client, err :=docker.NewClient(
			c.getDockerExternalClient(),
			*c.configs.Docker,
		)
		if err != nil {
			panic(fmt.Errorf("get docker client: %w", err))
		}

		return client
	})
}
