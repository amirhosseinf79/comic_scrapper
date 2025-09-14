package broker

import (
	"encoding/json"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
	"github.com/hibiken/asynq"
)

type clientS struct {
	client *asynq.Client
}

const (
	typePageProcess = "page:process"
)

func NewClient(addr, pwd string) interfaces.AsynqClient {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})
	return &clientS{client: client}
}

func (c *clientS) EnqueueTask(task *asynq.Task) (*asynq.TaskInfo, error) {
	return c.client.Enqueue(
		task,
		asynq.Timeout(30*time.Minute),
		asynq.MaxRetry(3),
	)
}

func (c *clientS) NewPageProcess(fields manager.PerPageScrap) (*asynq.Task, error) {
	payload, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(typePageProcess, payload), nil
}
