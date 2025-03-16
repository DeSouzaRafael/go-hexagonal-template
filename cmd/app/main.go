package app

import (
	"log"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/repositories"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/service"
)

func main() {

	config.InitConfig()

	db, err := database.NewDatabaseAdapter(config.AppConfig.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(&userService)

	server := web.NewWebService(userHandler)
	log.Fatal(server.Start())
}
