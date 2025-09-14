package interfaces

import "github.com/gofiber/fiber/v2"

type ScraperHandler interface {
	RequestProcess(ctx *fiber.Ctx) error
	GetLogByID(ctx *fiber.Ctx) error
}
