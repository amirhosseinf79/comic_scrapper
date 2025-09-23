package manager

import "github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"

type PerPageScrap struct {
	LogID uint   `json:"logId"`
	Page  string `json:"page"`
	WebhookRequest
}

type SendWebhook struct {
	WebhookRequest
	ComicInfo comic.Info
	LogID     uint
}
