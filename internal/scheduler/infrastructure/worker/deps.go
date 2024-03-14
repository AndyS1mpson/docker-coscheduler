package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"google.golang.org/grpc"
)

type externalClient interface {
	Build(ctx context.Context, in *task.BuildRequest, opts ...grpc.CallOption) (*task.BuildResponse, error)
}
