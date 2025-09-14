package server

func (s server) InitScrapHandlers() {
	user := s.app.Group("api/v1/scrapper")
	user.Post("/request", s.scrapperHandler.RequestProcess)
}
