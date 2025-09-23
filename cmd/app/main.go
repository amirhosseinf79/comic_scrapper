package main

import (
	"fmt"
	"os"

	"github.com/amirhosseinf79/comic_scrapper/cmd/job"
	scraprequest "github.com/amirhosseinf79/comic_scrapper/internal/application/handler/scrap_request"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/broker"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/database"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/persistence"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/server"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/logger"
	manager2 "github.com/amirhosseinf79/comic_scrapper/internal/service/manager"
	"github.com/lpernett/godotenv"

	_ "github.com/amirhosseinf79/comic_scrapper/docs"
)

// @title Comic Scrapper
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer [...]
// @schemes http
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	mainPort := os.Getenv("PORT")
	redisServer := os.Getenv("RedisServer")
	redisPassword := os.Getenv("RedisPass")
	dbConnStr := os.Getenv("SQLDB")
	debugStr := os.Getenv("DEBUG")
	debug := false
	if debugStr == "true" {
		debug = true
	}

	client := broker.NewClient(redisServer, redisPassword)
	queueHandler := broker.NewQueueServer(redisServer, redisPassword)
	db := database.NewGormConnection(dbConnStr, debug)
	logRepo := persistence.NewLoggerRepo(db)
	go job.HandleBroker(db, client, queueHandler)

	logService := logger.NewLoggerService(logRepo)
	manager := manager2.NewScrapperManager(client, logService)
	handler := scraprequest.NewManagerHandler(manager, logService)

	server1 := server.NewWebServer(handler)
	server1.InitSwaggerHandlers()
	server1.InitScrapHandlers()
	server1.InitLoggerHandlers()
	server1.Start(mainPort)

}
