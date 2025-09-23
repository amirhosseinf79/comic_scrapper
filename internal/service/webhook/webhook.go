package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/shared"
)

type webhook struct {
	client *http.Client
}

func NewWebhook() interfaces.WebhookService {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return &webhook{
		client: client,
	}
}

func (s *webhook) SendComicInfo(header manager.WebhookRequest, comicInfo comic.Info) error {
	jsonBody, err := json.Marshal(comicInfo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", header.WebhookURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", header.Authorization)
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return shared.ErrInvalidRequest
	}
	return nil
}
