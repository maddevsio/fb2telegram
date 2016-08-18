package conf

import (
	"os"

	"github.com/gen1us2k/log"
	"github.com/urfave/cli"
)

// Version stores current service version
var (
	Version              string
	FacebookClientID     string
	FacebookClientSecret string
	FacebookPageName     string
	TelegramBotToken     string
	TelegramWebhookURL   string
	HTTPBindAddr         string
	LogLevel             string
)

type Configuration struct {
	data *Fb2TelegramConfig
	app  *cli.App
}

// NewConfigurator is constructor and creates a new copy of Configuration
func NewConfigurator() *Configuration {
	Version = "0.1dev"
	app := cli.NewApp()
	app.Name = "Facebook page to telegram bot"
	app.Usage = "Show latest posts from facebook page"
	return &Configuration{
		data: &Fb2TelegramConfig{},
		app:  app,
	}
}

func (c *Configuration) fillConfig() *Fb2TelegramConfig {
	return &Fb2TelegramConfig{
		FacebookClientID:     FacebookClientID,
		FacebookClientSecret: FacebookClientSecret,
		FacebookPageName:     FacebookPageName,
		TelegramBotToken:     TelegramBotToken,
		TelegramWebhookURL:   TelegramWebhookURL,
		HTTPBindAddr:         HTTPBindAddr,
	}
}

// Run is wrapper around cli.App
func (c *Configuration) Run() error {
	c.app.Before = func(ctx *cli.Context) error {
		log.SetLevel(log.MustParseLevel(LogLevel))
		return nil
	}
	c.app.Flags = c.setupFlags()
	return c.app.Run(os.Args)
}

// App is public method for Configuration.app
func (c *Configuration) App() *cli.App {
	return c.app
}

func (c *Configuration) setupFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "http_bind_addr",
			Value:       ":8090",
			Usage:       "Set address to bind http server",
			EnvVar:      "HTTP_BIND_ADDR",
			Destination: &HTTPBindAddr,
		},
		cli.StringFlag{
			Name:        "facebook_client_id",
			Value:       "",
			Usage:       "Set facebook client id",
			EnvVar:      "FACEBOOK_CLIENT_ID",
			Destination: &FacebookClientID,
		},
		cli.StringFlag{
			Name:        "facebook_client_secret",
			Value:       "",
			Usage:       "set facebook client secret",
			EnvVar:      "FACEBOOK_CLIENT_SECRET",
			Destination: &FacebookClientSecret,
		},
		cli.StringFlag{
			Name:        "facebook_page_name",
			Value:       "ololohaus",
			Usage:       "Set page name from facebook",
			EnvVar:      "FACEBOOK_PAGE_NAME",
			Destination: &FacebookPageName,
		},
		cli.StringFlag{
			Name:        "telegram_bot_token",
			Value:       "",
			Usage:       "Set telegram bot access token",
			EnvVar:      "TELEGRAM_TOKEN",
			Destination: &TelegramBotToken,
		},
		cli.StringFlag{
			Name:        "telegram_web_hook_url",
			Value:       "",
			Usage:       "Set telegram bot webhook url",
			EnvVar:      "TELEGRAM_WEBHOOK_URL",
			Destination: &TelegramWebhookURL,
		},
		cli.StringFlag{
			Name:        "loglevel",
			Value:       "debug",
			Usage:       "set log level",
			Destination: &LogLevel,
			EnvVar:      "LOG_LEVEL",
		},
	}

}

// Get returns filled BillginConfig
func (c *Configuration) Get() *Fb2TelegramConfig {
	c.data = c.fillConfig()
	return c.data
}
