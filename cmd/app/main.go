package main

import (
	"log"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
)

func main() {

	// Load settings
	if err := config.LoadConfig(); err != nil {
		panic("Error loading settings: " + err.Error())
	}

	// Init database
	db, err := database.NewDatabaseAdapter(config.AppConfig.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		panic("Error migrating database: " + err.Error())
	}

	// Init container
	cont := container.NewContainer(db)

	// Init web service
	server := web.NewWebService(cont.Handlers)
	log.Fatal(server.Start())
}
