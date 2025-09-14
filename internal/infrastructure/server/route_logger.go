package server

func (s server) InitLoggerHandlers() {
	user := s.app.Group("api/v1/logger")
	user.Get("/:id", s.scrapperHandler.GetLogByID)
}
