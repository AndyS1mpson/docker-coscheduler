package models

var (
	// TaskStatusCreated новая задача в системе
	TaskStatusNew TaskStatus = "new"
	// TaskStatusBuild docker образ для задачи собран
	TaskStatusBuild TaskStatus = "build"
	// TaskStatusCreated задача создана, для нее создан docker контейнер
	TaskStatusCreated TaskStatus = "created"
)

// TaskStatus статус задачи
type TaskStatus string

// Task абстракция, описывающая задачу
type Task struct {
	ID      string     `json:"id"`      // Идентификатор задачи
	Node    *Node      `json:"node"`    // Нода, назначенная задаче для выполнения
	ImageID string     `json:"imageId"` // Идентификатор docker образа
	Title   string     `json:"title"`   // Описание задачи
	Status  TaskStatus `json:"status"`  // Статус задачи
	Config  *Config    `json:"config"`  // Конфигурация задачи
}

// Config конфигурация задачи
type Config struct {
	ContainerID string  `json:"containerId"`    // Идентификатор docker контейнера
	CPUs        *CPUSet `json:"cpus,omitempty"` // Конфигурация CPU
}
