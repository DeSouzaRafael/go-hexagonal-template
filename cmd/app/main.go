package main

import (
	"log"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/repositories"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/service"
)

func main() {

	if err := config.LoadConfig(); err != nil {
		panic("Erro ao carregar configurações: " + err.Error())
	}

	db, err := database.NewDatabaseAdapter(config.AppConfig.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&domain.User{})

	userRepository := repositories.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(&authService)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(&userService)

	server := web.NewWebService(userHandler, authHandler)
	log.Fatal(server.Start())
}
