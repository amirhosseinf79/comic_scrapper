package manager

type PerPageScrap struct {
	LogID         uint   `json:"logId"`
	Page          string `json:"page"`
	Authorization string `query:"authorization"`
	WebhookURL    string `json:"webhookUrl"`
}
