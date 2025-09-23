package interfaces

type AsynqServer interface {
	AddServices(client AsynqClient, webhook WebhookService, scrapper Scrapper, logger LoggerService) AsynqServer
	Start()
}
