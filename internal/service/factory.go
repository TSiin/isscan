package service

import (
	"isscan/pkg/config"
	"isscan/pkg/dns"
	"isscan/pkg/proxy"
	"net/http"

	"go.uber.org/zap"
)

// ServiceFactory 服务工厂接口
type ServiceFactory interface {
	NewPingService() PingService
	NewPortService() PortService
	NewSSLService() SSLService
	NewWebSocketService() WebSocketService
}

type serviceFactory struct {
	config   *config.Config
	logger   *zap.Logger
	client   *http.Client
	dialer   *proxy.Dialer
	resolver *dns.Resolver
}

func NewServiceFactory(cfg *config.Config) ServiceFactory {
	// 获取全局logger
	logger := zap.L()

	// 创建DNS解析器
	resolver := dns.NewResolver(cfg.DNS.Servers, cfg.DNS.Timeout)

	// 创建代理拨号器go
	dialer := proxy.NewDialer(&cfg.Proxy, resolver)

	// 创建HTTP传输层
	transport := dialer.GetTransport()

	// 创建HTTP客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Scanner.Timeout,
	}

	return &serviceFactory{
		config:   cfg,
		logger:   logger,
		client:   client,
		dialer:   dialer,
		resolver: resolver,
	}
}

func (f *serviceFactory) NewPingService() PingService {
	return NewPingService(f.config)
}

func (f *serviceFactory) NewPortService() PortService {
	return NewPortService(f.config, f.logger)
}

func (f *serviceFactory) NewSSLService() SSLService {
	return NewSSLService(f.config)
}

func (f *serviceFactory) NewWebSocketService() WebSocketService {
	return NewWebSocketService(f.config, f.logger)
}
