package server

import (
	"log"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/gofiber/fiber/v2"
)

type server struct {
	scrapperHandler interfaces.ScraperHandler
	app             *fiber.App
}

func NewWebServer(
	scrapperHandler interfaces.ScraperHandler,
) interfaces.ServerService {
	app := fiber.New()

	return &server{
		app:             app,
		scrapperHandler: scrapperHandler,
	}
}

func (s server) Start(port string) {
	err := s.app.Listen(":" + port)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
