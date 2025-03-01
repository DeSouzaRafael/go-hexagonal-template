package app

import (
	"log"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
)

func main() {

	config.InitConfig()

	db, err := database.NewConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := web.NewWebService()
	log.Fatal(server.Start())
}
