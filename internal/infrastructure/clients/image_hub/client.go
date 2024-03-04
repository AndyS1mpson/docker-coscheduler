package image_hub

// Client клиент для взаимодействие с хранилищем образов докер
type Client struct {
	workerImageHubDir string
}

// NewClient конструктор для Client
func NewClient(config Config) *Client {
	return &Client{
		workerImageHubDir: config.WorkerImageHubDir,
	}
}
