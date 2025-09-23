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
	WebhookError      string         `json:"webhookError"`
	HasInfo           bool           `json:"hasInfo"`
	TimeEstimated     int64          `json:"timeEstimated"`
	TotalEpisodes     int            `json:"totalEpisodes"`
	TotalFiles        int            `json:"totalFiles"`
	ProcessedEpisodes int            `json:"processedEpisodes"`
	ProcessedFiles    int            `json:"processedFiles"`
	Output            comic.Info     `json:"output" gorm:"serializer:json"`
	Console           []string       `json:"console" gorm:"serializer:json"`
}

func InitLog() *Log {
	return &Log{
		HasInfo: true,
		Status:  enum.Queued,
		Console: make([]string, 0),
	}
}

func (l *Log) SetStatus() {
	if l.ProcessedFiles == l.TotalFiles && l.ProcessedEpisodes == l.TotalEpisodes && l.TotalEpisodes != 0 && l.TotalFiles != 0 {
		l.Status = enum.Succeed
	} else {
		l.Status = enum.Failed
	}
}
