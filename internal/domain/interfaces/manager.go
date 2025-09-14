package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
)

type Manager interface {
	GenerateJson(fields manager.PageScrapRequest) ([]manager.PerPageScrap, error)
}
