package proxy

import (
	"context"
	"isscan/pkg/config"
	"isscan/pkg/dns"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

type Dialer struct {
	config   *config.ProxyConfig
	resolver *dns.Resolver
}

func NewDialer(cfg *config.ProxyConfig, resolver *dns.Resolver) *Dialer {
	return &Dialer{
		config:   cfg,
		resolver: resolver,
	}
}

func (d *Dialer) GetTransport() *http.Transport {
	transport := &http.Transport{
		DialContext: d.DialContext,
	}

	if d.config.Enable {
		proxyURL, err := url.Parse(d.config.URL)
		if err == nil {
			switch d.config.Type {
			case "http", "https":
				transport.Proxy = http.ProxyURL(proxyURL)
			case "socks5":
				auth := &proxy.Auth{
					User:     d.config.Username,
					Password: d.config.Password,
				}
				if dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct); err == nil {
					transport.DialContext = dialer.(proxy.ContextDialer).DialContext
				}
			}
		}
	}

	return transport
}

func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	// 使用DoH解析域名
	ips, err := d.resolver.LookupIP(ctx, host)
	if err != nil {
		return nil, err
	}

	// 尝试连接所有IP
	var lastErr error
	for _, ip := range ips {
		conn, err := net.DialTimeout(network, net.JoinHostPort(ip.String(), port), d.config.Timeout)
		if err == nil {
			return conn, nil
		}
		lastErr = err
	}

	return nil, lastErr
}
