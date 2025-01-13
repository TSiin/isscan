package app

import (
	"isscan/internal/service"
	"isscan/pkg/config"
)

type Application struct {
	Config           *config.Config
	ServiceFactory   service.ServiceFactory
	PingService      service.PingService
	PortService      service.PortService
	SSLService       service.SSLService
	WebSocketService service.WebSocketService
}

func NewApplication(cfg *config.Config) *Application {
	factory := service.NewServiceFactory(cfg)

	return &Application{
		Config:           cfg,
		ServiceFactory:   factory,
		PingService:      factory.NewPingService(),
		PortService:      factory.NewPortService(),
		SSLService:       factory.NewSSLService(),
		WebSocketService: factory.NewWebSocketService(),
	}
}
