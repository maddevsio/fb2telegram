package service

import (
	"fmt"
	"time"

	"github.com/gen1us2k/log"
	fb "github.com/huandu/facebook"
)

type FacebookService struct {
	BaseService

	logger  log.Logger
	fb2tg   *Fb2Telegram
	fbc     *fb.App
	fbToken string
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
	fbToken := fs.fbc.AppAccessToken()
	res, err := fb.Get(
		fmt.Sprintf("/%s/events", fs.fb2tg.Config().FacebookPageName),
		fb.Params{"access_token": fbToken},
	)
	if err != nil {
		return err
	}

	var items []fb.Result

	err = res.DecodeField("data", &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		startTime, err := time.Parse("2006-01-02T15:04:05-0700", item["start_time"].(string))
		if err != nil {
			fs.logger.Errorf("error thile parsing time: %s", err)
		}
		if startTime.Sub(time.Now()).Hours() >= 1 {
			fs.logger.Infof("%s post: %s\n", item["start_time"], item["description"])
		}
	}
	return nil
}
