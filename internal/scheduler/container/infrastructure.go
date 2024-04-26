package container

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/storage"
	imageHub "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/task_hub"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/container"
)

// GetStorage провайдер сконфигурированного хранилища
func (c *Container) GetStorage() *storage.Storage {
	return container.MustOrGetNew(c.Container, func() *storage.Storage {
		return storage.New(
			c.getPgConnection(),
		)
	})
}

// getPgConnection провайдер сконфигурированного подключения в БД
func (c *Container) getPgConnection() *sqlx.DB {
	return container.MustOrGetNew(c.Container, func() *sqlx.DB {
		conn, err := sqlx.Open(
			"postgres",
			fmt.Sprintf(
				"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				c.configs.Database.Host,
				c.configs.Database.Port,
				c.configs.Database.User,
				c.configs.Database.Password,
				c.configs.Database.DB,
			))
		if err != nil {
			panic(fmt.Sprintf("init storage conn: %s", err))
		}

		conn.SetMaxOpenConns(10)
		conn.SetMaxIdleConns(5)
		conn.SetConnMaxLifetime(time.Minute * 15)

		c.PushShutdown(func() {
			_ = conn.Close()
		})

		return conn
	})
}

func (c *Container) GetTaskHub() *imageHub.Hub {
	return container.MustOrGetNew(c.Container, func() *imageHub.Hub {
		return imageHub.NewHub(
			c.configs.TaskDir,
		)
	})
}
