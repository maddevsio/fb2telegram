package conf

// Fb2TelegramConfig stores all configuration of service
type Fb2TelegramConfig struct {
	FacebookClientID     string
	FacebookClientSecret string
	FacebookPageName     string
	TelegramBotToken     string
	TelegramWebhookURL   string
	HTTPBindAddr         string
}
