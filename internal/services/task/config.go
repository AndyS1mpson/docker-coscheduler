package task

// Config конфигурация воркера
type Config struct {
	CPUNums int64  // Колисество ядер на ноде
	URI     string `envconfig:"SERVICE_HOST" required:"true"` // URI ноды
	Port    int64  `envconfig:"SERVICE_PORT" required:"true"` // Порт ноды
}
