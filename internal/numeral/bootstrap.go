package numeral

import (
	"github.com/erecarte/showcase/internal/numeral/api"
	"github.com/erecarte/showcase/internal/numeral/api/payment_orders"
	"github.com/erecarte/showcase/internal/numeral/database"
	"log"
)

type App struct {
	server *api.HttpServer
	config *Config
}

func NewApp(cfg *Config) *App {
	return &App{
		config: cfg,
	}
}

func (app *App) Start() {
	db, err := database.GetConnection(app.config.DbLocation)
	if err != nil {
		log.Fatal(err)
	}
	storage := payment_orders.NewSQLLiteStorage(db)
	fileGenerator := payment_orders.NewBankFileSender(app.config.BankFolderLocation)
	service := payment_orders.NewService(storage, fileGenerator)
	fileReciever, _ := payment_orders.NewBankFileReceiver(app.config.BankFolderLocation, service)
	err = fileReciever.ReceiveFilesFromBank()
	if err != nil {
		log.Fatal(err)
	}

	app.server = api.NewHttpServer(app.config.Port, service)
	app.server.Start()
}

func (app *App) Stop() {
	app.server.Stop()
}
