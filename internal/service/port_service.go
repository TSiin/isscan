package service

import (
	"context"
	"isscan/pkg/config"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type portService struct {
	*BaseService
	concurrent int
	timeout    time.Duration
}

func NewPortService(cfg *config.Config, logger *zap.Logger) PortService {
	return &portService{
		BaseService: NewBaseService(cfg, logger),
		concurrent:  cfg.Scanner.Concurrent,
		timeout:     cfg.Scanner.Timeout,
	}
}

func (s *portService) ScanPort(ctx context.Context, host string, port int) (*PortResult, error) {
	result := &PortResult{
		Port:    port,
		Service: s.detectService(port),
	}

	start := time.Now()
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		if !strings.Contains(err.Error(), "refused") {
			s.logger.Debug("Port scan error",
				zap.String("host", host),
				zap.Int("port", port),
				zap.Error(err))
		}
		result.Error = err.Error()
		return result, nil
	}
	defer conn.Close()

	result.Open = true
	result.Latency = time.Since(start)

	return result, nil
}

func (s *portService) ScanPorts(ctx context.Context, host string, ports []int) ([]*PortResult, error) {
	// 如果没有指定端口，使用默认端口列表
	if len(ports) == 0 {
		s.logger.Debug("Using default ports for scanning",
			zap.Ints("ports", s.config.Scanner.DefaultPorts))
		ports = s.config.Scanner.DefaultPorts
	}

	results := make([]*PortResult, len(ports))
	var wg sync.WaitGroup
	sem := make(chan struct{}, s.concurrent)

	for i, port := range ports {
		wg.Add(1)
		go func(i, port int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// 设置单个端口扫描的超时
			scanCtx, cancel := context.WithTimeout(ctx, s.timeout)
			defer cancel()

			result, err := s.ScanPort(scanCtx, host, port)
			if err != nil {
				s.logger.Error("Port scan failed",
					zap.String("host", host),
					zap.Int("port", port),
					zap.Error(err))
			}
			results[i] = result
		}(i, port)
	}

	wg.Wait()

	// 过滤掉关闭的端口
	openPorts := make([]*PortResult, 0)
	for _, result := range results {
		if result != nil && result.Open {
			openPorts = append(openPorts, result)
		}
	}

	s.logger.Info("Port scan completed",
		zap.String("host", host),
		zap.Int("total_ports", len(ports)),
		zap.Int("open_ports", len(openPorts)))

	return openPorts, nil
}

func (s *portService) detectService(port int) string {
	// 常见端口服务映射
	services := map[int]string{
		80:    "HTTP",
		443:   "HTTPS",
		22:    "SSH",
		21:    "FTP",
		3306:  "MySQL",
		5432:  "PostgreSQL",
		6379:  "Redis",
		27017: "MongoDB",
	}

	if service, ok := services[port]; ok {
		return service
	}
	return ""
}
