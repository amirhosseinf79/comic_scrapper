package model

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Status            enum.LogStatus `json:"status"`
	WebhookSend       bool           `json:"webhookSend"`
	HasInfo           bool           `json:"hasInfo"`
	TotalEpisodes     int            `json:"totalEpisodes"`
	TotalFiles        int            `json:"totalFiles"`
	ProcessedEpisodes int            `json:"processedEpisodes"`
	ProcessedFiles    int            `json:"processedFiles"`
	Console           []string       `json:"console" gorm:"serializer:json"`
	Output            comic.Info     `json:"-" gorm:"serializer:json"`
}

func InitLog() *Log {
	return &Log{
		Status:  enum.Queued,
		HasInfo: true,
		Console: make([]string, 0),
	}
}
