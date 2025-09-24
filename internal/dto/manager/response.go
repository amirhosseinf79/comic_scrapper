package manager

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
)

type PerPageResponse struct {
	LogID uint `json:"logId"`
}

type LogMock struct {
	ID                uint           `json:"ID"`
	Status            enum.LogStatus `json:"status"`
	HasInfo           bool           `json:"hasInfo"`
	TimeEstimated     int64          `json:"timeEstimated"`
	TotalEpisodes     int            `json:"totalEpisodes"`
	TotalFiles        int            `json:"totalFiles"`
	ProcessedEpisodes int            `json:"processedEpisodes"`
	ProcessedFiles    int            `json:"processedFiles"`
	Console           []string       `json:"console"`
	Output            comic.Info     `json:"output"`
}
