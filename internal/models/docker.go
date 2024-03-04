package models

var (
	// ContainerStateCreated контейнер создан
	ContainerStateCreated ContainerState = "created"
	// ContainerStateRunning контейнер запущен
	ContainerStateRunning ContainerState = "running"
	// ContainerStatePaused контейнер остановлен
	ContainerStatePaused ContainerState = "paused"
	// SateRestarting контейнер перезапускается
	SateRestarting ContainerState = "restarting"
	// ContainerStateRemoving контейнер удаляется
	ContainerStateRemoving ContainerState = "removing"
	// ContainerStateExited контейнер завершен
	ContainerStateExited ContainerState = "exited"
	// ContainerStateDead контейнер убит
	ContainerStateDead ContainerState = "dead"
)

// ContainerState состояние контейнера
type ContainerState string

// ContainerInfo информация о контейнере
type ContainerInfo struct {
	ID       string         // Идентификатор контейнера
	State    ContainerState // Состояние контейнера
	ExitCode int64          // Статус код контейнера
}
