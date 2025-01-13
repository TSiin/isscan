package ping

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type PingResult struct {
	Target     string        `json:"target"`
	IPv4       string        `json:"ipv4,omitempty"`
	IPv6       string        `json:"ipv6,omitempty"`
	MinLatency time.Duration `json:"minLatency"`
	MaxLatency time.Duration `json:"maxLatency"`
	AvgLatency time.Duration `json:"avgLatency"`
	PacketLoss float64       `json:"packetLoss"`
	SendCount  int           `json:"sendCount"`
	RecvCount  int           `json:"recvCount"`
	Protocol   string        `json:"protocol"`
	Error      string        `json:"error,omitempty"`
}

type Pinger struct {
	Count    int
	Interval time.Duration
	Timeout  time.Duration
}

func NewPinger() *Pinger {
	return &Pinger{
		Count:    3,
		Interval: time.Second,
		Timeout:  time.Second * 5,
	}
}

func (p *Pinger) Ping(target string) PingResult {
	result := PingResult{
		Target:    target,
		SendCount: p.Count,
	}

	// 解析IP地址
	ips, err := net.LookupIP(target)
	if err != nil {
		result.Error = fmt.Sprintf("解析域名失败: %v", err)
		return result
	}

	// 分离IPv4和IPv6地址
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			result.IPv4 = ipv4.String()
		} else if ip.To16() != nil {
			result.IPv6 = ip.String()
		}
	}

	if result.IPv4 == "" && result.IPv6 == "" {
		result.Error = "未找到有效的IP地址"
		return result
	}

	// 优先使用IPv4
	pingAddr := target
	if result.IPv4 != "" {
		result.Protocol = "IPv4"
		pingAddr = result.IPv4
	} else {
		result.Protocol = "IPv6"
		pingAddr = result.IPv6
	}

	// 根据操作系统选择ping命令
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", strconv.Itoa(p.Count), pingAddr)
	case "linux", "darwin":
		cmd = exec.Command("ping", "-c", strconv.Itoa(p.Count), pingAddr)
	default:
		result.Error = "不支持的操作系统"
		return result
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Sprintf("ping执行失败: %v", err)
		return result
	}

	// 解析ping结果
	outputStr := string(output)

	// 解析延迟时间
	var totalMs float64
	var minMs, maxMs float64
	var received int

	// 根据操作系统解析输出
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") || strings.Contains(line, "时间=") {
			received++
			parts := strings.Split(line, "time=")
			if len(parts) == 2 {
				ms := parseFloat(parts[1])
				totalMs += ms
				if minMs == 0 || ms < minMs {
					minMs = ms
				}
				if ms > maxMs {
					maxMs = ms
				}
			}
		}
	}

	result.RecvCount = received
	if received > 0 {
		result.MinLatency = time.Duration(minMs * float64(time.Millisecond))
		result.MaxLatency = time.Duration(maxMs * float64(time.Millisecond))
		result.AvgLatency = time.Duration((totalMs / float64(received)) * float64(time.Millisecond))
		result.PacketLoss = float64(p.Count-received) / float64(p.Count) * 100
	} else {
		result.Error = "所有ping请求超时"
	}

	return result
}

// 辅助函数：解析延迟时间字符串
func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "ms")
	s = strings.TrimSpace(s)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
