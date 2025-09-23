package persistence

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/repository"
	"gorm.io/gorm"
)

type logS struct {
	db *gorm.DB
}

func NewLoggerRepo(db *gorm.DB) repository.Logger {
	return &logS{db: db}
}

func (l logS) Create(log *model.Log) error {
	return l.db.Create(log).Error
}

func (l logS) Update(log *model.Log) error {
	return l.db.Save(log).Error
}

func (l logS) GetById(id uint) (*model.Log, error) {
	var log model.Log
	err := l.db.First(&log, id).Error
	return &log, err
}

func (l logS) GetListById(ids []uint) ([]model.Log, error) {
	var logs []model.Log
	if len(ids) == 0 {
		return logs, gorm.ErrRecordNotFound
	}
	err := l.db.Where("id IN ? and status = ?", ids, enum.Succeed).Find(&logs).Error
	return logs, err
}
