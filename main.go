package main

import (
	"context"
	"flag"
	"fmt"
	"isscan/internal/app"
	"isscan/internal/middleware"
	"isscan/internal/router"
	"isscan/pkg/config"
	"isscan/pkg/logger"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/zap"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.InitLogger(config.GlobalConfig.Log); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d",
			config.GlobalConfig.Server.Host,
			config.GlobalConfig.Server.Port,
		)),
	)

	application := app.NewApplication(&config.GlobalConfig)

	h.Use(middleware.Recovery())
	h.Use(middleware.RequestLogger())

	router.Register(h, application.ServiceFactory)

	logger.Info("Server starting",
		zap.String("host", config.GlobalConfig.Server.Host),
		zap.Int("port", config.GlobalConfig.Server.Port),
	)

	// 添加版本信息日志
	logger.Info("Starting isscan",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("build_date", date),
	)

	// 创建根context，不设置超时
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down gracefully...")
		cancel()
	}()

	// 启动服务器
	go func() {
		if err := h.Run(); err != nil {
			logger.Error("Error starting server", zap.Error(err))
			cancel()
		}
	}()

	// 等待退出信号
	<-ctx.Done()

	// 设置关闭超时时间
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// 优雅关闭服务器
	if err := h.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during shutdown", zap.Error(err))
	}

	// 等待所有请求处理完成
	time.Sleep(5 * time.Second)

	// 关闭其他资源
	if err := zap.L().Sync(); err != nil {
		fmt.Printf("Error syncing logger: %v\n", err)
	}
}
