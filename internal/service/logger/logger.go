package logger

import (
	"fmt"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/repository"
)

type loggerService struct {
	baseRepo repository.Logger
}

func NewLoggerService(baseRepo repository.Logger) interfaces.LoggerService {
	return &loggerService{
		baseRepo: baseRepo,
	}
}

func (l *loggerService) Create(log *model.Log) error {
	return l.baseRepo.Create(log)
}

func (l *loggerService) AutoUpdate(log *model.Log, state string, status enum.LogStatus, cmd ...string) error {
	msg := fmt.Sprintf("%v: %v - %v:", state, status.String(), cmd)
	log.Console = append(log.Console, msg)
	fmt.Println(msg)
	//if status == enum.Failed {
	//	log.Status = status
	//}
	if len(log.Console)%5 == 0 {
		return l.Update(log)
	}
	return nil
}

func (l *loggerService) Update(log *model.Log) error {
	return l.baseRepo.Update(log)
}

func (l *loggerService) GetById(id uint) (*model.Log, error) {
	logM, err := l.baseRepo.GetById(id)
	return logM, err
}

func (l *loggerService) GetListById(ids []uint) ([]model.Log, error) {
	logM, err := l.baseRepo.GetListById(ids)
	return logM, err
}
