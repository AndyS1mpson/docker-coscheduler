package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type externalClient interface {
	BuildTask(ctx context.Context, in *task.BuildTaskRequest, opts ...grpc.CallOption) (*task.BuildTaskResponse, error)
	CreateTask(ctx context.Context, in *task.CreateTaskRequest, opts ...grpc.CallOption) (*task.CreateTaskResponse, error)
	PauseTask(ctx context.Context, in *task.PauseTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
