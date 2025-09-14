package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
	"github.com/hibiken/asynq"
)

type serverS struct {
	server   *asynq.Server
	scrapper interfaces.Scrapper
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
	scrapper interfaces.Scrapper,
	logger interfaces.LoggerService,
) interfaces.AsynqServer {
	s.scrapper = scrapper
	s.logger = logger
	return s
}

func (s *serverS) HandlePageProcess(ctx context.Context, t *asynq.Task) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var p manager.PerPageScrap
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Processing %v, logID: %v\n", p.Page, p.LogID)
	logM, err := s.logger.GetById(p.LogID)
	if err != nil {
		return err
	}

	c := make(chan error, 1)
	go func() {
		_, err := s.scrapper.GenerateComicInfo(logM, p.Page)
		if logM.ProcessedFiles == logM.TotalFiles && logM.ProcessedEpisodes == logM.TotalEpisodes {
			logM.Status = enum.Succeed
		} else {
			logM.Status = enum.Failed
		}
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
		return res
	}
}

func (s *serverS) Start() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(typePageProcess, s.HandlePageProcess)

	if err := s.server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
