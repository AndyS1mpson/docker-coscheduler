package container

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/worker"
)

func (c *Container) GetWorkerClient(uri string, port int64) *worker.Client {
	url := fmt.Sprintf("%s:%d", uri, port)

	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("connect to worker client: %w", err))
	}

	extClient := task.NewTaskClient(conn)

	return worker.NewClient(extClient)
}
