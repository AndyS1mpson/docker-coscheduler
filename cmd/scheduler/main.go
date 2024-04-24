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

	appContainer, shutdown := container.NewContainer(*config)
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
		client := appContainer.GetWorkerClient(node.Host, node.Port)
		info, err := client.GetNodeInfo(appContainer.Ctx())
		if err != nil {
			log.Error(err, log.Data{"host": node.Host, "port": node.Port})

			return failExitCode
		}

		nodeClients[*info] = client
	}

	strategies := map[string]container.Strategy{
		"roundRobin": appContainer.GetRoundRobinStrategy(nodeClients),
		"fcs":        appContainer.GetFCSStrategy(nodeClients),
		"fcn":        appContainer.GetFCNStrategy(nodeClients),
	}

	for name, strategy := range strategies {
		duration, err := strategy.Execute(appContainer.Ctx(), tasks)
		if err != nil {
			log.Error(fmt.Errorf("%s strategy: %w", name, err), log.Data{})

			return failExitCode
		}

		log.Info(fmt.Sprintf("%s execution total time: %v", name, duration), log.Data{})
	}

	return successExitCode
}
