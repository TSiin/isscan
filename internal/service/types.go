package service

import (
	"context"
	"time"
)

// Service interfaces
type PingService interface {
	Ping(ctx context.Context, host string) (*PingResult, error)
}

type PortService interface {
	ScanPort(ctx context.Context, host string, port int) (*PortResult, error)
	ScanPorts(ctx context.Context, host string, ports []int) ([]*PortResult, error)
}

type SSLService interface {
	CheckDomain(ctx context.Context, domain string) (*SSLResult, error)
	BatchCheck(ctx context.Context, domains []string) (map[string]*SSLResult, error)
}

type WebSocketService interface {
	CheckEndpoint(ctx context.Context, host, path string, secure bool) (*WebSocketResult, error)
	BatchCheck(ctx context.Context, endpoints []Endpoint) ([]*WebSocketResult, error)
}

// Result types
type PingResult struct {
	Target     string        `json:"target"`
	IPv4       string        `json:"ipv4,omitempty"`
	IPv6       string        `json:"ipv6,omitempty"`
	MinLatency time.Duration `json:"min_latency"`
	MaxLatency time.Duration `json:"max_latency"`
	AvgLatency time.Duration `json:"avg_latency"`
	PacketLoss float64       `json:"packet_loss"`
	SendCount  int           `json:"send_count"`
	RecvCount  int           `json:"recv_count"`
	Protocol   string        `json:"protocol"`
}

type PortResult struct {
	Port    int           `json:"port"`
	Open    bool          `json:"open"`
	Service string        `json:"service,omitempty"`
	Error   string        `json:"error,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}

type SSLResult struct {
	Valid         bool      `json:"valid"`
	NotBefore     time.Time `json:"not_before"`     // 证书生效时间
	NotAfter      time.Time `json:"not_after"`      // 证书过期时间
	DaysRemaining int       `json:"days_remaining"` // 剩余有效天数
	DNSNames      []string  `json:"dns_names"`
	Issuer        string    `json:"issuer"`
	Subject       string    `json:"subject"`       // 证书主体
	SerialNumber  string    `json:"serial_number"` // 序列号
	Error         string    `json:"error,omitempty"`
}

type WebSocketResult struct {
	Host      string        `json:"host"`
	Path      string        `json:"path"`
	Protocol  string        `json:"protocol"`
	Available bool          `json:"available"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
}

type Endpoint struct {
	Host   string `json:"host"`
	Path   string `json:"path"`
	Secure bool   `json:"secure"`
}
