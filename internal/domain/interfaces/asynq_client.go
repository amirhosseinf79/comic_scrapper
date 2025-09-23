package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
	"github.com/hibiken/asynq"
)

type AsynqClient interface {
	NewPageProcess(fields manager.PerPageScrap) (*asynq.Task, error)
	NewWebhookSend(fields manager.SendWebhook) (*asynq.Task, error)
	EnqueueTask(task *asynq.Task) (*asynq.TaskInfo, error)
}
