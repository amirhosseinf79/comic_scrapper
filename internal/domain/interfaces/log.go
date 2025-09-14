package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
)

type LoggerService interface {
	Create(log *model.Log) error
	AutoUpdate(log *model.Log, state string, status enum.LogStatus, cmd ...string) error
	Update(log *model.Log) error
	GetById(id uint) (*model.Log, error)
}
