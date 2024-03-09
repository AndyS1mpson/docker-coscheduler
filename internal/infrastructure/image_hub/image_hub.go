package image_hub

// Hub хранилищем образов докер
type Hub struct {
	workerImageHubDir string
}

// NewHub конструктор для Hub
func NewHub(config Config) *Hub {
	return &Hub{
		workerImageHubDir: config.WorkerImageHubDir,
	}
}
