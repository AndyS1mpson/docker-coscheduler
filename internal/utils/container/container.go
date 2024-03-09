package container

import (
	"context"
	"sync"

	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/datastructs"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/shutdown"
)

// ServiceName ключ сервиса для резолвинга уже созданного инстанса зависимости
type ServiceName string

// Container DI-контейнер сервиса
type Container struct {
	appCtx context.Context

	services sync.Map
	shutdown *datastructs.Stack[func()]
}

// NewContainer консутрктор для Container
func NewContainer() (c *Container, gracefulShutdown func()) {
	ctx, cancel := shutdown.WithCancel(context.Background())

	c = &Container{
		appCtx:   ctx,
		shutdown: &datastructs.Stack[func()]{},
	}

	gracefulShutdown = c.gracefulShutdown(cancel)

	return
}

// Ctx получение контекста выполнения приложения
func (c *Container) Ctx() context.Context {
	return c.appCtx
}

// PushShutdown добавить функцию завершения в стек
func (c *Container) PushShutdown(f func()) {
	c.shutdown.Push(f)
}

func (c *Container) gracefulShutdown(ctxCancel func()) func() {
	return func() {
		for !c.shutdown.IsEmpty() {
			c.shutdown.Pop()()
		}

		ctxCancel()
	}
}
