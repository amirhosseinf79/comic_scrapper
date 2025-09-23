package job

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/infrastructure/persistence"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/logger"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/scrapper"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/webhook"
	"gorm.io/gorm"
)

func HandleBroker(db *gorm.DB, client interfaces.AsynqClient, handler interfaces.AsynqServer) {
	logRepo := persistence.NewLoggerRepo(db)
	logService := logger.NewLoggerService(logRepo)
	scrapper2 := scrapper.New(true, logService)
	defer scrapper2.Close()
	webhookService := webhook.NewWebhook()
	handler.AddServices(client, webhookService, scrapper2, logService)
	handler.Start()
}
