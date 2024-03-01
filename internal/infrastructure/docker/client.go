package docker

import "context"

// Client клиент для работы с Docker Engine API
type Client struct {
	externalClient externalClient
}

// NewClient конструктор для Client
func NewClient(client externalClient) (*Client, error) {
	return &Client{
		externalClient: client,
	}, nil
}

// Ping проверка работоспособности docker engine на сервере
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.externalClient.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
