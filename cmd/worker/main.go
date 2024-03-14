package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"github.com/AndyS1mpson/docker-coscheduler/internal/worker/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/worker/controllers"
	"github.com/AndyS1mpson/docker-coscheduler/internal/worker/infrastructure/middleware"
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

	// Create TCP connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", container.GetConfigs().Node.Port))
	if err != nil {
		log.Error(fmt.Errorf("failed to listen: %w", err), log.Data{})

		return failExitCode
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.MiddlewareLog,
		),
	)
	reflection.Register(s)

	task.RegisterTaskServer(s, controllers.NewServer(container.GetTaskService()))

	if err = s.Serve(lis); err != nil {
		log.Error(fmt.Errorf("failed to server: %w", err), log.Data{})

		return exitCode
	}

	return successExitCode
}
