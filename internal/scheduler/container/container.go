package container

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/dotenv"
)

var configPath = ""

// Container DI-контейнер сервиса
type Container struct {
	configs AppConfig
	*container.Container
}

// NewContainer конструктор для Container
func NewContainer(config AppConfig) (c *Container, gracefulShutdown func()) {
	dotenv.Load()

	sc, gracefulShutdown := container.NewContainer()

	c = &Container{
		configs:   config,
		Container: sc,
	}

	return
}

// GetConfigs получение заполненной конфигурации
func (c *Container) GetConfigs() AppConfig {
	return c.configs
}
