package image_hub

// Config конфигурация клиента
type Config struct {
	WorkerImageHubDir string `envconfig:"WORKER_IMAGE_HUB_DIR" required:"true"`
}
