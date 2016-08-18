package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/gen1us2k/fb2telegram/conf"
	"github.com/gen1us2k/fb2telegram/service"
	"github.com/gen1us2k/log"
	"github.com/urfave/cli"
)

func main() {
	app := conf.NewConfigurator()
	app.App().Action = func(ctx *cli.Context) error {
		worker := service.NewFb2Telegram(app.Get())
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		defer signal.Stop(signalChan)

		go func() {
			<-signalChan
			log.Info("signal received, stopping...")
			worker.Stop()

			time.Sleep(2 * time.Second)
			os.Exit(0)
		}()

		err := worker.Start()
		if err != nil {
			log.Fatalf("error on local node start, %v", err)
		}

		worker.WaitStop()
		return nil
	}
	if err := app.Run(); err != nil {
		log.Fatalf("Error on run app, %v", err)
	}
}
