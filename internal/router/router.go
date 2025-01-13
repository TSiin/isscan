package router

import (
	"isscan/internal/handler"
	"isscan/internal/service"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(h *server.Hertz, factory service.ServiceFactory) {
	// 创建统一的handler
	handlers := handler.NewHandler(factory)

	// API分组
	api := h.Group("/api")
	{
		// 基础检查
		api.POST("/ping", handlers.PingHandler)

		// 端口扫描
		port := api.Group("/port")
		{
			port.POST("/scan", handlers.ScanPortHandler)
		}

		// SSL证书检查
		ssl := api.Group("/ssl")
		{
			ssl.POST("/check", handlers.CheckSSLHandler)
			ssl.POST("/batch", handlers.BatchCheckSSLHandler)
		}

		// WebSocket检测
		ws := api.Group("/websocket")
		{
			ws.POST("/check", handlers.CheckWSHandler)
			ws.POST("/batch", handlers.BatchCheckWSHandler)
		}
	}
}
