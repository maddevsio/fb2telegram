package service

import (
	"fmt"

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
		fmt.Sprintf("/%s/posts", fs.fb2tg.Config().FacebookPageName),
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
		if _, ok := item["message"]; ok {
			fs.logger.Infof("%s post: %s\n", item["created_time"], item["message"])
		}
	}
	return nil
}
