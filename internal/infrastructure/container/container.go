package container

import (
	"github.com/mrfansi/go-api-boilerplate/internal/application/service"
	domainRepository "github.com/mrfansi/go-api-boilerplate/internal/domain/repository"
	"github.com/mrfansi/go-api-boilerplate/internal/infrastructure/config"
	"github.com/mrfansi/go-api-boilerplate/internal/infrastructure/database"
	infraRepository "github.com/mrfansi/go-api-boilerplate/internal/infrastructure/repository"
	"github.com/mrfansi/go-api-boilerplate/internal/interfaces/http/handler"
	"github.com/mrfansi/go-api-boilerplate/internal/interfaces/http/middleware"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Container struct {
	container *dig.Container
}

func NewContainer() *Container {
	return &Container{
		container: dig.New(),
	}
}

func (c *Container) Configure(cfg *config.Config) error {
	// Provide config
	if err := c.container.Provide(func() *config.Config {
		return cfg
	}); err != nil {
		return err
	}

	// Provide database
	if err := c.container.Provide(database.NewSQLiteDB); err != nil {
		return err
	}

	// Provide repositories
	if err := c.container.Provide(func(db *gorm.DB) domainRepository.UserRepository {
		return infraRepository.NewUserRepository(db)
	}); err != nil {
		return err
	}

	// Provide services
	if err := c.container.Provide(service.NewAuthService); err != nil {
		return err
	}
	if err := c.container.Provide(service.NewUserService); err != nil {
		return err
	}

	// Provide handlers
	if err := c.container.Provide(handler.NewAuthHandler); err != nil {
		return err
	}
	if err := c.container.Provide(handler.NewUserHandler); err != nil {
		return err
	}

	// Provide middleware
	if err := c.container.Provide(middleware.NewAuthMiddleware); err != nil {
		return err
	}
	if err := c.container.Provide(middleware.NewLoggerMiddleware); err != nil {
		return err
	}
	if err := c.container.Provide(middleware.NewCorsMiddleware); err != nil {
		return err
	}

	return nil
}

func (c *Container) Resolve(constructor interface{}) error {
	return c.container.Invoke(constructor)
}

func (c *Container) Container() *dig.Container {
	return c.container
}
