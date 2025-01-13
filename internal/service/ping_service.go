package service

import (
	"context"
	"isscan/pkg/config"
	"time"

	"github.com/go-ping/ping"
)

type pingService struct {
	config  *config.Config
	timeout time.Duration
	count   int
}

func NewPingService(cfg *config.Config) PingService {
	return &pingService{
		config:  cfg,
		timeout: cfg.Scanner.Timeout,
		count:   3,
	}
}

func (s *pingService) Ping(ctx context.Context, host string) (*PingResult, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, err
	}

	pinger.Count = s.count
	pinger.Timeout = s.timeout
	pinger.SetPrivileged(false)

	err = pinger.Run()
	if err != nil {
		return nil, err
	}

	stats := pinger.Statistics()
	return &PingResult{
		Target:     host,
		MinLatency: stats.MinRtt,
		MaxLatency: stats.MaxRtt,
		AvgLatency: stats.AvgRtt,
		PacketLoss: stats.PacketLoss,
		SendCount:  stats.PacketsSent,
		RecvCount:  stats.PacketsRecv,
	}, nil
}
