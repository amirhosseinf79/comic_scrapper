package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/shared"
	"github.com/amirhosseinf79/comic_scrapper/internal/service/scrapper"
	"github.com/hibiken/asynq"
)

type serverS struct {
	client   interfaces.AsynqClient
	server   *asynq.Server
	scrapper interfaces.Scrapper
	webhook  interfaces.WebhookService
	logger   interfaces.LoggerService
}

func NewQueueServer(addr, pwd string) interfaces.AsynqServer {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     addr,
			Password: pwd,
			DB:       0,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return 2 * time.Second
			},
		},
	)
	return &serverS{server: server}
}

func (s *serverS) AddServices(
	client interfaces.AsynqClient,
	webhook interfaces.WebhookService,
	scrapper interfaces.Scrapper,
	logger interfaces.LoggerService,
) interfaces.AsynqServer {
	s.client = client
	s.webhook = webhook
	s.scrapper = scrapper
	s.logger = logger
	return s
}

func (s *serverS) Start() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(typePageProcess, s.handlePageProcess)
	mux.HandleFunc(typeSendWebhook, s.handleSendWebhook)

	if err := s.server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (s *serverS) handlePageProcess(ctx context.Context, t *asynq.Task) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var p manager.PerPageScrap
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	nowTime := time.Now().Unix()
	log.Printf("Processing %v, logID: %v\n", p.Page, p.LogID)
	logM, err := s.logger.GetById(p.LogID)
	logM.SetPendingStatus()
	if err != nil {
		return err
	}

	c := make(chan error, 1)
	go func() {
		_, err := scrapper.GenerateComicInfo(s.scrapper, logM, p.Page)
		logM.SetFinalStatus()
		logM.TimeEstimated = time.Now().Unix() - nowTime
		_ = s.logger.Update(logM)
		select {
		case c <- err:
		case <-ctx.Done():
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-c:
		if res != nil {
			return res
		}
		logM, err := s.logger.GetById(p.LogID)
		if err != nil {
			return err
		}
		req := manager.SendWebhook{
			WebhookRequest: p.WebhookRequest,
			ComicInfo:      logM.Output,
			LogID:          p.LogID,
		}
		task, err := s.client.NewWebhookSend(req)
		if err != nil {
			return err
		}
		_, err = s.client.EnqueueTask(task)
		return err
	}
}

func (s *serverS) handleSendWebhook(_ context.Context, t *asynq.Task) error {
	var p manager.SendWebhook
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logM, err := s.logger.GetById(p.LogID)
	if err != nil {
		return err
	}
	err = s.webhook.SendComicInfo(p.WebhookRequest, p.ComicInfo)
	if err != nil {
		_ = s.logger.AutoUpdate(logM, "SendWebhook", enum.Failed, p.WebhookURL, err.Error())
		if errors.Is(err, shared.ErrInvalidRequest) {
			log.Printf("Invalid webhook request: %v", err.Error())
			err = fmt.Errorf("invalid webhook request: %v %w", err.Error(), asynq.SkipRetry)
		}
	} else {
		_ = s.logger.AutoUpdate(logM, "SendWebhook", enum.Succeed, p.WebhookURL)
	}
	_ = s.logger.Update(logM)
	return err
}
