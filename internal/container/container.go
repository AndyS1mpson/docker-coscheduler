package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/dotenv"
)

// Container DI-контейнер сервиса
type Container struct {
	configs AppConfig
	*container.Container
}

// NewContainer конструктор для Container
func NewContainer() (c *Container, gracefulShutdown func()) {
	dotenv.Load()

	sc, gracefulShutdown := container.NewContainer()

	c = &Container{
		configs:   LoadConfig(),
		Container: sc,
	}

	return
}

// GetConfigs получение заполненной конфигурации
func (c *Container) GetConfigs() AppConfig {
	return c.configs
}
