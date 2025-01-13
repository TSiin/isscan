package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"isscan/pkg/config"
	"net"
	"sync"
	"time"
)

type sslService struct {
	config  *config.Config
	timeout time.Duration
}

func NewSSLService(cfg *config.Config) SSLService {
	return &sslService{
		config:  cfg,
		timeout: cfg.SSL.Timeout,
	}
}

func (s *sslService) CheckDomain(ctx context.Context, domain string) (*SSLResult, error) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: s.timeout},
		"tcp",
		net.JoinHostPort(domain, "443"),
		&tls.Config{ServerName: domain},
	)
	if err != nil {
		return &SSLResult{
			Valid: false,
			Error: err.Error(),
		}, nil
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	now := time.Now()

	return &SSLResult{
		Valid:         now.After(cert.NotBefore) && now.Before(cert.NotAfter),
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		DaysRemaining: int(cert.NotAfter.Sub(now).Hours() / 24),
		DNSNames:      cert.DNSNames,
		Issuer:        cert.Issuer.CommonName,
		Subject:       cert.Subject.CommonName,
		SerialNumber:  cert.SerialNumber.String(),
	}, nil
}

func (s *sslService) BatchCheck(ctx context.Context, domains []string) (map[string]*SSLResult, error) {
	results := make(map[string]*SSLResult)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()

			result, err := s.CheckDomain(ctx, domain)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("domain %s check failed: %w", domain, err):
				default:
				}
				return
			}

			mu.Lock()
			results[domain] = result
			mu.Unlock()
		}(domain)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return nil, err
	default:
		return results, nil
	}
}

// 实现 SSLService 接口方法
