package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
)

type ManagerService interface {
	GenerateJson(fields manager.PageScrapRequest) ([]manager.PerPageResponse, error)
	SendWebhookByLogID(fields manager.SendWebhookRequest) ([]manager.PerPageResponse, error)
}
