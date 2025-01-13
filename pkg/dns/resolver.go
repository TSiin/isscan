package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Resolver struct {
	servers []string
	client  *http.Client
	timeout time.Duration
}

func NewResolver(servers []string, timeout time.Duration) *Resolver {
	return &Resolver{
		servers: servers,
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (r *Resolver) LookupIP(ctx context.Context, host string) ([]net.IP, error) {
	for _, server := range r.servers {
		ips, err := r.queryDoH(ctx, server, host)
		if err == nil && len(ips) > 0 {
			return ips, nil
		}
	}
	// 如果DoH失败，回退到系统DNS
	return net.DefaultResolver.LookupIP(ctx, "ip4", host)
}

func (r *Resolver) queryDoH(ctx context.Context, server, host string) ([]net.IP, error) {
	url := fmt.Sprintf("%s?name=%s&type=A", server, host)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/dns-json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Answer []struct {
			Data string `json:"data"`
		} `json:"Answer"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, ans := range result.Answer {
		if ip := net.ParseIP(ans.Data); ip != nil {
			ips = append(ips, ip)
		}
	}

	return ips, nil
}
