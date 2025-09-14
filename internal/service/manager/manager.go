package manager

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
)

type managerS struct {
	asynqClient interfaces.AsynqClient
	logger      interfaces.LoggerService
}

func NewScrapperManager(
	asynqClient interfaces.AsynqClient,
	logger interfaces.LoggerService,
) interfaces.ManagerService {
	return &managerS{
		asynqClient: asynqClient,
		logger:      logger,
	}
}

func (m managerS) GenerateJson(fields manager.PageScrapRequest) ([]manager.PerPageResponse, error) {
	finalLogs := make([]manager.PerPageResponse, 0)
	for _, page := range fields.Pages {
		log := model.InitLog()
		if err := m.logger.Create(log); err != nil {
			return nil, err
		}
		field := manager.PerPageScrap{
			Authorization: fields.Authorization,
			WebhookURL:    fields.WebhookURL,
			Page:          page,
			LogID:         log.ID,
		}
		finalLogs = append(finalLogs, manager.PerPageResponse{LogID: log.ID})
		process, err := m.asynqClient.NewPageProcess(field)
		if err != nil {
			return nil, err
		}
		_, err = m.asynqClient.EnqueueTask(process)
		if err != nil {
			return nil, err
		}
	}
	return finalLogs, nil
}
