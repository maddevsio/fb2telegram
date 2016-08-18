package service

import (
	"fmt"
	"sync"

	"github.com/gen1us2k/fb2telegram/conf"
	"github.com/gen1us2k/log"
)

// Fb2Telegram is main struct of daemon
// it stores all services that used by
type Fb2Telegram struct {
	config *conf.Fb2TelegramConfig

	services  map[string]Service
	waitGroup sync.WaitGroup

	logger log.Logger
}

// NewFb2Telegram creates and returns new Fb2TelegramInstance
func NewFb2Telegram(config *conf.Fb2TelegramConfig) *Fb2Telegram {
	pb := new(Fb2Telegram)
	pb.config = config
	pb.logger = log.NewLogger("fb2telegram")
	pb.services = make(map[string]Service)
	pb.AddService(&FacebookService{})
	pb.AddService(&TelegramService{})
	return pb
}

// Start starts all services in separate goroutine
func (pb *Fb2Telegram) Start() error {
	pb.logger.Info("Starting bot service")
	for _, service := range pb.services {
		pb.logger.Infof("Initializing: %s\n", service.Name())
		if err := service.Init(pb); err != nil {
			return fmt.Errorf("initialization of %q finished with error: %v", service.Name(), err)
		}
		pb.waitGroup.Add(1)

		go func(srv Service) {
			defer pb.waitGroup.Done()
			pb.logger.Infof("running %q service\n", srv.Name())
			if err := srv.Run(); err != nil {
				pb.logger.Errorf("error on run %q service, %v", srv.Name(), err)
			}
		}(service)
	}
	return nil
}

// AddService adds service into Fb2Telegram.services map
func (pb *Fb2Telegram) AddService(srv Service) {
	pb.services[srv.Name()] = srv

}

// Config returns current instance of Fb2TelegramConfig
func (pb *Fb2Telegram) Config() conf.Fb2TelegramConfig {
	return *pb.config
}

// Stop stops all services running
func (pb *Fb2Telegram) Stop() {
	pb.logger.Info("Worker is stopping...")
	for _, service := range pb.services {
		service.Stop()
	}
}

// WaitStop blocks main thread and waits when all goroutines will be stopped
func (pb *Fb2Telegram) WaitStop() {
	pb.waitGroup.Wait()
}

func (pb *Fb2Telegram) FacebookService() *FacebookService {
	service, ok := pb.services["facebook_service"]
	if !ok {
		pb.logger.Error("Error getting facebook_service")
	}
	return service.(*FacebookService)
}
