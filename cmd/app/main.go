// @title           Go Hexagonal Template API
// @version         1.0
// @description     RESTful API template built with Hexagonal Architecture, Echo, and GORM.
// @host            localhost:8086
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"flag"
	"log"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/DeSouzaRafael/go-hexagonal-template/pkg/util"

	_ "github.com/DeSouzaRafael/go-hexagonal-template/docs"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mockDB := flag.Bool("mock-db", false, "Use mock database instead of real database")
	flag.Parse()

	if err := config.LoadConfig(); err != nil {
		return err
	}

	var db port.Database
	var err error

	if *mockDB {
		db, err = database.NewMockDatabaseAdapter()
	} else {
		db, err = database.NewDatabaseAdapter(config.AppConfig.Database)
	}

	if err != nil {
		return err
	}
	defer db.Close()

	return runWithDB(db)
}

type WebServer interface {
	Start() error
}

func runWithDB(db port.Database) error {
	webServiceFactory := func(h container.Handlers) WebServer {
		return web.NewWebService(h)
	}
	return runWithDependencies(db, container.NewContainer, webServiceFactory)
}

type ContainerFactory func(port.Database) *container.Container

type WebServiceFactory func(container.Handlers) WebServer

func runWithDependencies(db port.Database, containerFactory ContainerFactory, webServiceFactory WebServiceFactory) error {
	if !util.CurrentExecutionEnvironmentProduction() {
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			return err
		}
	}

	cont := containerFactory(db)

	server := webServiceFactory(cont.Handlers)
	return server.Start()
}
