package service

import (
	"isscan/pkg/config"
	"net/url"
	"time"

	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"crypto/tls"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type wsService struct {
	config  *config.Config
	timeout time.Duration
	logger  *zap.Logger
}

func NewWebSocketService(cfg *config.Config, logger *zap.Logger) WebSocketService {
	return &wsService{
		config:  cfg,
		timeout: cfg.WebSocket.Timeout,
		logger:  logger,
	}
}

func (s *wsService) CheckEndpoint(ctx context.Context, host, path string, secure bool) (*WebSocketResult, error) {
	scheme := "ws"
	if secure {
		scheme = "wss"
	}

	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	result := &WebSocketResult{
		Host:      host,
		Path:      path,
		Protocol:  scheme,
		Available: false,
	}

	// 配置 WebSocket dialer
	dialer := websocket.Dialer{
		HandshakeTimeout: s.timeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 可选：忽略证书验证
		},
		Proxy: http.ProxyFromEnvironment,
	}

	// 添加重试逻辑
	var lastErr error
	for i := 0; i < s.config.WebSocket.RetryTimes; i++ {
		start := time.Now()

		// 先检查TCP连接
		_, err := net.DialTimeout("tcp", net.JoinHostPort(host, s.getPort(secure)), s.timeout)
		if err != nil {
			lastErr = fmt.Errorf("TCP connection failed: %v", err)
			s.logger.Debug("TCP connection attempt failed",
				zap.String("host", host),
				zap.Int("attempt", i+1),
				zap.Error(err))
			time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
			continue
		}

		// 尝试WebSocket连接
		c, resp, err := dialer.DialContext(ctx, u.String(), nil)
		if err != nil {
			if resp != nil {
				lastErr = fmt.Errorf("WebSocket connection failed with status %d: %v", resp.StatusCode, err)
			} else {
				lastErr = fmt.Errorf("WebSocket connection failed: %v", err)
			}
			s.logger.Debug("WebSocket connection attempt failed",
				zap.String("url", u.String()),
				zap.Int("attempt", i+1),
				zap.Error(err))
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// 连接成功
		defer c.Close()
		result.Available = true
		result.Latency = time.Since(start)
		return result, nil
	}

	// 所有重试都失败
	result.Error = lastErr.Error()
	return result, nil
}

func (s *wsService) getPort(secure bool) string {
	if secure {
		return "443"
	}
	return "80"
}

func (s *wsService) BatchCheck(ctx context.Context, endpoints []Endpoint) ([]*WebSocketResult, error) {
	results := make([]*WebSocketResult, len(endpoints))
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, endpoint := range endpoints {
		wg.Add(1)
		go func(i int, endpoint Endpoint) {
			defer wg.Done()

			result, err := s.CheckEndpoint(ctx, endpoint.Host, endpoint.Path, endpoint.Secure)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}

			mu.Lock()
			results[i] = result
			mu.Unlock()
		}(i, endpoint)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return nil, err
	default:
		return results, nil
	}
}

// 实现 WebSocketService 接口方法
