package controllers

import "github.com/AndyS1mpson/docker-coscheduler/generated/task"

// Server реализация gRPC сервера
type Server struct {
	task.UnimplementedTaskServer
	service service
}

// NewServer конструктор для Server
func NewServer(service service) *Server {
	return &Server{
		service: service,
	}
}
