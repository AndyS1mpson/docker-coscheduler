package image_hub

// Hub хранилищем образов докер
type Hub struct {
	schedulerTaskDir string
}

// NewHub конструктор для Hub
func NewHub(taskHubDir string) *Hub {
	return &Hub{
		schedulerTaskDir: taskHubDir,
	}
}
