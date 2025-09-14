package interfaces

import (
	"context"

	"github.com/hibiken/asynq"
)

type AsynqServer interface {
	AddServices(scrapper Scrapper, logger LoggerService) AsynqServer
	HandlePageProcess(ctx context.Context, t *asynq.Task) error
	Start()
}
