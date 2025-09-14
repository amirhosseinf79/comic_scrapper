package model

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Status            enum.LogStatus `json:"status"`
	HasInfo           bool           `json:"hasInfo"`
	TotalEpisodes     int            `json:"totalEpisodes"`
	TotalFiles        int            `json:"totalFiles"`
	ProcessedEpisodes int            `json:"processedEpisodes"`
	ProcessedFiles    int            `json:"processedFiles"`
	Console           []string       `json:"-"`
	Output            comic.Info     `json:"output" gorm:"serializer:json"`
}

func InitLog() *Log {
	return &Log{
		Status:  enum.Queued,
		HasInfo: true,
	}
}
