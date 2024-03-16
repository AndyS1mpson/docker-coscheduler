package worker

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

type externalClient interface {
	BuildTask(ctx context.Context, in *task.BuildTaskRequest, opts ...grpc.CallOption) (*task.BuildTaskResponse, error)
	CreateTask(ctx context.Context, in *task.CreateTaskRequest, opts ...grpc.CallOption) (*task.CreateTaskResponse, error)
	StartTask(ctx context.Context, in *task.StartTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PauseTask(ctx context.Context, in *task.PauseTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ResumeTask(ctx context.Context, in *task.ResumeTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	StopTask(ctx context.Context, in *task.StopTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
