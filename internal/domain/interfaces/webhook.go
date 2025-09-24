package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
)

type WebhookService interface {
	SendComicInfo(header manager.WebhookRequest, comicInfo comic.Info) error
	SetLog(log *model.Log, logger LoggerService) WebhookService
}
