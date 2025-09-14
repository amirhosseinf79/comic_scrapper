package repository

import "github.com/amirhosseinf79/comic_scrapper/internal/domain/model"

type Logger interface {
	Create(log *model.Log) error
	Update(log *model.Log) error
	GetById(id uint) (*model.Log, error)
}
