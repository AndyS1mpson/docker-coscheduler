package image_hub

// Hub хранилищем образов докер
type Hub struct {
	imageHubDir string
}

// NewHub конструктор для Hub
func NewHub(dir string) *Hub {
	return &Hub{
		imageHubDir: dir,
	}
}
