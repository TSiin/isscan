package websocket

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type WSResult struct {
	Host      string        `json:"host"`
	Path      string        `json:"path"`
	Protocol  string        `json:"protocol"`  // ws 或 wss
	Available bool          `json:"available"` // 是否可用
	Latency   time.Duration `json:"latency"`   // 连接延迟
	Error     string        `json:"error,omitempty"`
	Response  string        `json:"response,omitempty"`
}

type WSChecker struct {
	Timeout    time.Duration
	RetryTimes int
	RetryDelay time.Duration
	SkipVerify bool
}

func NewWSChecker() *WSChecker {
	return &WSChecker{
		Timeout:    10 * time.Second,
		RetryTimes: 3,
		RetryDelay: time.Second,
	}
}

func (w *WSChecker) CheckWithContext(ctx context.Context, host, path string, secure bool) WSResult {
	result := WSResult{
		Host:     cleanHost(host),
		Path:     ensurePath(path),
		Protocol: getProtocol(secure),
	}

	client := &http.Client{
		Timeout: w.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: w.SkipVerify,
			},
			DisableKeepAlives: true,
		},
	}

	for i := 0; i < w.RetryTimes; i++ {
		select {
		case <-ctx.Done():
			result.Error = "context canceled"
			return result
		default:
			if success := w.tryConnect(client, &result); success {
				return result
			}
			if i < w.RetryTimes-1 {
				time.Sleep(w.RetryDelay)
			}
		}
	}
	return result
}

func (w *WSChecker) tryConnect(client *http.Client, result *WSResult) bool {
	url := buildURL(result.Host, result.Path, result.Protocol == "wss")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = fmt.Sprintf("create request failed: %v", err)
		return false
	}

	addWSHeaders(req)

	start := time.Now()
	resp, err := client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Error = fmt.Sprintf("request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusSwitchingProtocols {
		result.Available = true
		return true
	}

	if resp.StatusCode == http.StatusBadRequest {
		result.Response = "handshake error: bad \"Upgrade\" header"
		result.Available = true // WebSocket端点存在，但需要正确的握手
		return true
	}

	result.Response = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
	return false
}

// 辅助函数
func cleanHost(host string) string {
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "ws://")
	return strings.TrimPrefix(host, "wss://")
}

func ensurePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func getProtocol(secure bool) string {
	if secure {
		return "wss"
	}
	return "ws"
}

func buildURL(host, path string, secure bool) string {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func addWSHeaders(req *http.Request) {
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Sec-WebSocket-Version", "13")
	req.Header.Add("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
}
