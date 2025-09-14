package server

import fiberSwagger "github.com/swaggo/fiber-swagger"

func (s server) InitSwaggerHandlers() {
	s.app.Get("/swagger/*", fiberSwagger.WrapHandler)
}
