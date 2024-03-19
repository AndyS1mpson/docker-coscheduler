package main

import (
	"fmt"
	"os"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
)

const (
	successExitCode = 0
	failExitCode    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	var err error

	config, err := container.NewConfig()
	if err != nil {
		log.Error(fmt.Errorf("read config: %w", err), log.Data{})

		return failExitCode
	}

	container, shutdown := container.NewContainer(*config)

	defer func() {
		if panicErr := recover(); panicErr != nil {
			exitCode = failExitCode
		}

		if err != nil {
			exitCode = failExitCode
		}
	}()

	defer shutdown()

	tasks := []models.StrategyTask{
		{
			Name:       "task 1",
			FolderName: "task1",
		},
		// {
		// 	Name: "task 2",
		// 	FolderName: "task1",
		// },
	}

	extClient := container.GetWorkerClient(config.Nodes[0].Host, config.Nodes[0].Port)

	strategy := container.GetSequentialStrategy([]*worker.Client{extClient})

	strategy.Execute(container.Ctx(), tasks)

	return successExitCode
}
