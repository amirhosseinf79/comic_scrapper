package interfaces

type ServerService interface {
	InitLoggerHandlers()
	InitScrapHandlers()
	Start(port string)
}
