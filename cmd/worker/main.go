package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/container"
	"github.com/AndyS1mpson/docker-coscheduler/internal/controllers"
	"github.com/AndyS1mpson/docker-coscheduler/internal/infrastructure/middleware"
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

	container, shutdown := container.NewContainer()

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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", container.GetConfigs().NodeConfig.Port))
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
