package container

import (
	"runtime"

	"github.com/kelseyhightower/envconfig"

	"github.com/AndyS1mpson/docker-coscheduler/internal/infrastructure/clients/docker"
	imageHub "github.com/AndyS1mpson/docker-coscheduler/internal/infrastructure/image_hub"
	"github.com/AndyS1mpson/docker-coscheduler/internal/services/task"
)

// AppConfig структура, содержащая конфигурации менеджеров
type AppConfig struct {
	Docker     *docker.Config
	ImageHub   *imageHub.Config
	NodeConfig *task.Config
}

func LoadConfig() AppConfig {
	var appConfig AppConfig

	envconfig.MustProcess("", &appConfig)

	appConfig.NodeConfig.CPUNums = int64(runtime.NumCPU())

	return appConfig
}
