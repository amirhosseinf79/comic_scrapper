package scraprequest

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/manager"
	"github.com/gofiber/fiber/v2"
)

type scrapperH struct {
	manager interfaces.ManagerService
	logger  interfaces.LoggerService
}

func NewManagerHandler(manager interfaces.ManagerService, logger interfaces.LoggerService) interfaces.ScraperHandler {
	return &scrapperH{
		manager: manager,
		logger:  logger,
	}
}

// RequestProcess
// @Summary Request Scrapper
// @Description Request Scrapper to begin the process
// @Tags Service
// @Accept json
// @Produce json
// @Param fields body manager.PageScrapRequest true "Fields"
// @Success 200 {array} manager.PerPageResponse
// @Router /api/v1/scrapper/request [post]
func (s scrapperH) RequestProcess(ctx *fiber.Ctx) error {
	var fields manager.PageScrapRequest
	err := ctx.BodyParser(&fields)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	response, err := s.manager.GenerateJson(fields)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}

// GetLogByID
// @Summary Get Request Result
// @Description Get Request Result By LogID
// @Tags Log
// @Accept json
// @Produce json
// / @Param logID path int true "logID"
// @Success 200 {object} manager.LogMock
// @Router /api/v1/logger/{logID} [get]
func (s scrapperH) GetLogByID(ctx *fiber.Ctx) error {
	logID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	response, err := s.logger.GetById(uint(logID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}
