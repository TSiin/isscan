package utils

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ParseDuration 解析时间字符串
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// ValidateIPAddress 验证IP地址
func ValidateIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// GetHostIP 获取主机IP地址
func GetHostIP(host string) (string, string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", "", err
	}

	var ipv4, ipv6 string
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4 = ip.String()
		} else {
			ipv6 = ip.String()
		}
	}

	return ipv4, ipv6, nil
}

// FormatBytes 格式化字节大小
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// RetryWithTimeout 带超时的重试函数
func RetryWithTimeout(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := fn(); err == nil {
				return nil
			}
			if i < attempts-1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(sleep):
				}
			}
		}
	}
	return fmt.Errorf("达到最大重试次数")
}
