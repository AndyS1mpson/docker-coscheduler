package worker

// Client клиент для работы с воркерами
type Client struct {
	externalClient externalClient
}

// NewClient конструктор для Client
func NewClient(externalClient externalClient) *Client {
	return &Client{
		externalClient: externalClient,
	}
}
