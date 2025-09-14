package interfaces

type ServerService interface {
	InitLoggerHandlers()
	InitScrapHandlers()
	InitSwaggerHandlers()
	Start(port string)
}
