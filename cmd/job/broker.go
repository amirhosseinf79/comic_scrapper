package job

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/broker"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/persistence"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/logger"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/scrapper"
	"gorm.io/gorm"
)

func HandleBroker(db *gorm.DB) {
	//client := broker.NewClient("127.0.0.1", "6379", "")
	logRepo := persistence.NewLoggerRepo(db)
	logService := logger.NewLoggerService(logRepo)

	scrapper2 := scrapper.New(false, logService)
	defer scrapper2.Close()

	handler := broker.NewQueueServer("localhost", "6379", "")
	handler.AddServices(scrapper2, logService)
	handler.Start()
}
