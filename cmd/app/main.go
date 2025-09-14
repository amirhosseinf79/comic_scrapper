package main

import (
	"github.com/amirhosseinf79/comic_scrapper/cmd/job"
	scraprequest "github.com/amirhosseinf79/comic_scrapper/internal/application/handler/scrap_request"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/broker"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/database"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/persistence"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/server"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/logger"
	manager2 "github.com/amirhosseinf79/comic_scrapper/internal/service/manager"
)

func main() {
	client := broker.NewClient("localhost", "6379", "")
	db := database.NewGormConnection("", true)
	go job.HandleBroker(db)

	logRepo := persistence.NewLoggerRepo(db)
	logService := logger.NewLoggerService(logRepo)

	manager := manager2.NewScrapperManager(client, logService)

	handler := scraprequest.NewManagerHandler(manager, logService)

	server1 := server.NewWebServer(handler)
	server1.InitScrapHandlers()
	server1.InitLoggerHandlers()
	server1.Start("8080")

}
