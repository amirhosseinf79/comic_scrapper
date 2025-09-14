package manager

type PageScrapRequest struct {
	Authorization string   `query:"authorization"`
	WebhookURL    string   `json:"webhookUrl"`
	Pages         []string `json:"pages"`
}
