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

	tasks := make([]models.StrategyTask, 0, len(config.Tasks))

	for _, task := range config.Tasks {
		tasks = append(tasks, models.StrategyTask{
			Name:       task.Name,
			FolderName: task.FolderName,
		})
	}

	nodeClients := make(map[models.Node]*worker.Client)

	for _, node := range config.Nodes {
		client := container.GetWorkerClient(node.Host, node.Port)
		info, err := client.GetNodeInfo(container.Ctx())
		if err != nil {
			log.Error(err, log.Data{"host": node.Host, "port": node.Port})

			return failExitCode
		}

		nodeClients[*info] = client
	}

	// seqStrategy := container.GetSequentialStrategy(nodeClients)

	fcsStrategy := container.GetFCSStrategy(nodeClients)

	// seqDuration, err := seqStrategy.Execute(container.Ctx(), tasks)
	// if err != nil {
	// 	log.Error(fmt.Errorf("seq strategy: %w", err), log.Data{})

	// 	return failExitCode
	// }

	fcsDuration, err := fcsStrategy.Execute(container.Ctx(), tasks)
	if err != nil {
		log.Error(fmt.Errorf("fcs strategy: %w", err), log.Data{})

		return failExitCode
	}

	//log.Info(fmt.Sprintf("seq strategy execution total time: %v\n", seqDuration), log.Data{})
	log.Info(fmt.Sprintf("fcs strategy execution total time: %v", fcsDuration), log.Data{})

	return successExitCode
}
