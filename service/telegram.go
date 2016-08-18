package service

import (
	"fmt"
	"strings"

	"github.com/gen1us2k/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"gopkg.in/telegram-bot-api.v4"
)

type TelegramService struct {
	BaseService

	fb2tg        *Fb2Telegram
	logger       log.Logger
	e            *echo.Echo
	bot          *tgbotapi.BotAPI
	updateChan   chan tgbotapi.Update
	botCommands  map[string]func(*tgbotapi.Message)
	botChatTexts map[string]func(*tgbotapi.Message)
}

func (ts *TelegramService) Name() string {
	return "telegram_service"
}

func (ts *TelegramService) Init(fb2tg *Fb2Telegram) error {
	ts.fb2tg = fb2tg
	ts.logger = log.NewLogger(ts.Name())
	bot, err := tgbotapi.NewBotAPI(ts.fb2tg.Config().TelegramBotToken)
	if err != nil {
		return err
	}
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(ts.fb2tg.Config().TelegramWebhookURL))
	if err != nil {
		return err
	}

	ts.bot = bot
	ts.updateChan = make(chan tgbotapi.Update, 100)
	ts.e = echo.New()
	ts.e.POST("/", ts.handleBotRequests)
	ts.e.GET(fmt.Sprintf("/%s", bot.Token), ts.handleBotRequests)
	ts.e.POST(fmt.Sprintf("/%s", bot.Token), ts.handleBotRequests)
	ts.botCommands = make(map[string]func(*tgbotapi.Message))
	ts.botChatTexts = make(map[string]func(*tgbotapi.Message))
	ts.botCommands["/start"] = ts.onStartCommand
	ts.botCommands["/help"] = ts.onHelpCommand
	ts.botChatTexts["привет"] = ts.onStartCommand

	return nil
}

func (ts *TelegramService) Run() error {
	ts.fb2tg.waitGroup.Add(1)
	go ts.e.Run(standard.New(ts.fb2tg.Config().HTTPBindAddr))
	ts.handleRun()

	return nil
}

func (ts *TelegramService) handleBotRequests(c echo.Context) error {
	var update tgbotapi.Update
	if err := c.Bind(&update); err != nil {
		return err
	}
	ts.updateChan <- update
	return nil
}

func (ts *TelegramService) handleRun() {
	for update := range ts.updateChan {
		ts.logger.Infof("%+v\n", update)

		if update.Message == nil {
			continue
		}

		if update.Message.Chat.IsGroup() {
			_, err := ts.bot.GetMe()
			if err != nil {
				ts.logger.Errorf("Error: %s", err.Error())
				continue
			}
		}
		me, err := ts.bot.GetMe()
		if err != nil {
			ts.logger.Error("Error getting myself: %s", err)
		}
		if me.UserName != update.Message.From.UserName {
			message := strings.ToLower(update.Message.Text)
			ts.logger.Info(message)
			if _, ok := ts.botChatTexts[message]; ok {
				ts.botChatTexts[message](update.Message)
			}
		}

		if update.Message.IsCommand() {
			ts.logger.Info(update.Message.Text)
			if _, ok := ts.botCommands[update.Message.Text]; ok {
				ts.botCommands[update.Message.Text](update.Message)
			}
		}
	}
}

func (ts *TelegramService) onStartCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		`
            Привет! Я бот Ололохауса и я могу следующие вещи:
            /start - начать работу со мной
            /help - помощь
            /nearest - ближайшие наши события
            /remindme - напоминать тебе о событиях утром каждого дня, когда у нас будет событие
        `)
	msg.ReplyToMessageID = message.MessageID
	ts.bot.Send(msg)
}
func (ts *TelegramService) onHelpCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		`
            /start - начать работу со мной
            /help - помощь
            /nearest - ближайшие наши события
            /remindme - напоминать тебе о событиях утром каждого дня, когда у нас будет событие
        `)
	msg.ReplyToMessageID = message.MessageID
	ts.bot.Send(msg)
}
