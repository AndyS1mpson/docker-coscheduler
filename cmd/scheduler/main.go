package main

import (
	"fmt"
	"os"

	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/container"
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

	taskArchiver := container.GetTaskHub()

	archive, err := taskArchiver.ArchiveImageToTar("task1", "test_1")
	if err != nil {
		log.Error(err, log.Data{})

		return failExitCode
	}

	extClient := container.GetWorkerClient(config.Nodes[0].Host, config.Nodes[0].Port)

	imageID, err := extClient.Build(container.Ctx(), *archive, "test task")
	if err != nil {
		log.Error(err, log.Data{})

		return failExitCode
	}

	log.Println(imageID, log.Data{})

	return successExitCode
}
