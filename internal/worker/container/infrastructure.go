package container

import (
	"fmt"

	"github.com/docker/docker/client"

	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
	imageHub "github.com/AndyS1mpson/docker-coscheduler/internal/worker/infrastructure/image_hub"
)

// getDockerExternalClient получение внешнего клиента для взаимодействия с docker api
func (c *Container) getDockerExternalClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(fmt.Sprintf("init docker api engine client: %s", err))
	}

	c.PushShutdown(func() {
		cli.Close()
	})

	return cli
}

func (c *Container) GetImageHub() *imageHub.Hub {
	return container.MustOrGetNew(c.Container, func() *imageHub.Hub {
		return imageHub.NewHub(
			c.configs.WorkerImageHubDir,
		)
	})
}
