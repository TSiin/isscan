package service

import (
	"isscan/pkg/config"
	"time"

	"go.uber.org/zap"
)

type BaseService struct {
	config  *config.Config
	logger  *zap.Logger
	timeout time.Duration
}

func NewBaseService(cfg *config.Config, logger *zap.Logger) *BaseService {
	return &BaseService{
		config:  cfg,
		logger:  logger,
		timeout: cfg.Scanner.Timeout,
	}
}
