package shutdown

import (
	"context"
	"os/signal"
	"syscall"
)

// WithCancel возвращает контекст и функцию для отмены которые будут срабатывать на сигналы OS
func WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
}
