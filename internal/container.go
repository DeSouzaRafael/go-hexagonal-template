package container

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/repositories"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/service"
)

type Container struct {
	Repositories Repositories
	Handlers     Handlers
	Services     Services
}

type Repositories struct {
	UserRepository port.UserRepository
}

type Services struct {
	AuthService service.AuthService
	UserService service.UserService
}

type Handlers struct {
	AuthHandler handler.AuthHandler
	UserHandler handler.UserHandler
}

func NewContainer(db port.Database) *Container {
	// Init repositories
	userRepository := repositories.NewUserRepository(db)

	// Init services
	authService := service.NewAuthService(userRepository)
	userService := service.NewUserService(userRepository)

	// Init handlers
	authHandler := handler.NewAuthHandler(&authService)
	userHandler := handler.NewUserHandler(&userService)

	return &Container{
		Repositories: Repositories{
			UserRepository: userRepository,
		},
		Services: Services{
			AuthService: authService,
			UserService: userService,
		},
		Handlers: Handlers{
			AuthHandler: authHandler,
			UserHandler: userHandler,
		},
	}
}
