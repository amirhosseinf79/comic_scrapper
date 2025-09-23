package manager

type WebhookRequest struct {
	Authorization string `query:"authorization"`
	WebhookURL    string `json:"webhookUrl"`
}

type PageScrapRequest struct {
	WebhookRequest
	Pages []string `json:"pages"`
}

type SendWebhookRequest struct {
	WebhookRequest
	LogIDs []uint `json:"logIds"`
}
