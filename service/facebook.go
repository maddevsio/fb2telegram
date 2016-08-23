package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/gen1us2k/log"
	fb "github.com/huandu/facebook"
)

type FacebookService struct {
	BaseService

	logger log.Logger
	fb2tg  *Fb2Telegram
	fbc    *fb.App
	events []fb.Result
}

type FbItems struct {
	Data []fb.Result
}

func (fs *FacebookService) Name() string {
	return "facebook_service"
}
func (fs *FacebookService) Init(fb2tg *Fb2Telegram) error {
	fs.fb2tg = fb2tg
	fs.logger = log.NewLogger(fs.Name())
	fs.fbc = fb.New(
		fs.fb2tg.Config().FacebookClientID,
		fs.fb2tg.Config().FacebookClientSecret,
	)
	return nil
}
func (fs *FacebookService) Run() error {
	fs.updateEvents()
	for range time.Tick(time.Duration(3600) * time.Second) {
		fs.updateEvents()
	}
	return nil
}

func (fs *FacebookService) updateEvents() {
	fbToken := fs.fbc.AppAccessToken()
	res, err := fb.Get(
		fmt.Sprintf("/%s/events", fs.fb2tg.Config().FacebookPageName),
		fb.Params{"access_token": fbToken},
	)
	if err != nil {
		fs.logger.Errorf("Error getting token: %s", err)
		return
	}

	var data []fb.Result
	var events FbItems
	err = res.DecodeField("data", &data)
	if err != nil {
		fs.logger.Errorf("Error decoding Data: %s", err)
		return
	}
	fbItem := FbItems{
		Data: data,
	}
	for _, item := range fbItem.Data {
		startTime, err := time.Parse("2006-01-02T15:04:05-0700", item["start_time"].(string))
		if err != nil {
			fs.logger.Errorf("error thile parsing time: %s", err)
		}
		if startTime.Sub(time.Now()).Hours() >= 1 {
			events.Data = append(events.Data, item)
			fs.logger.Infof("%s post: %s\n", item["start_time"], item["description"])
		}
	}
	sort.Sort(events)
	fs.events = events.Data
}

func (fs *FacebookService) GetEventMessage() string {
	var message string
	for _, item := range fs.events {
		startTime, err := time.Parse("2006-01-02T15:04:05-0700", item["start_time"].(string))
		if err != nil {
			fs.logger.Errorf("Error while parsing date: %s", err)
		}
		message += startTime.Format("02-01 в 15:04 у нас ")

		message += item["name"].(string)
		message += "\n"
	}
	if message == "" {
		message = "Пока-что ближайших событий нет. Следите за обновлением"
	}
	return message
}

func (d FbItems) Len() int {
	return len(d.Data)
}

func (d FbItems) Less(i, j int) bool {
	return d.Data[i]["start_time"].(string) < d.Data[j]["start_time"].(string)
}

func (d FbItems) Swap(i, j int) {
	d.Data[i], d.Data[j] = d.Data[j], d.Data[i]
}
