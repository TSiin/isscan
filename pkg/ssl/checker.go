package ssl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

type CertInfo struct {
	Subject        string    `json:"subject"`        // 证书主体
	Issuer         string    `json:"issuer"`         // 颁发者
	SerialNumber   string    `json:"serialNumber"`   // 序列号
	Version        int       `json:"version"`        // 版本
	IsCA           bool      `json:"isCA"`           // 是否为 CA 证书
	NotBefore      time.Time `json:"notBefore"`      // 生效时间
	NotAfter       time.Time `json:"notAfter"`       // 过期时间
	SignatureAlg   string    `json:"signatureAlg"`   // 签名算法
	PublicKeyAlg   string    `json:"publicKeyAlg"`   // 公钥算法
	KeyUsage       []string  `json:"keyUsage"`       // 密钥用途
	DNSNames       []string  `json:"dnsNames"`       // DNS名称
	EmailAddresses []string  `json:"emailAddresses"` // 邮箱地址
}

type CertResult struct {
	Domain        string     `json:"domain"`
	ExpiryDate    time.Time  `json:"expiryDate"`
	RemainingDays int        `json:"remainingDays"`
	Status        bool       `json:"status"`
	Error         string     `json:"error,omitempty"`
	CertChain     []CertInfo `json:"certChain"`   // 证书链
	Protocol      string     `json:"protocol"`    // TLS 协议版本
	CipherSuite   string     `json:"cipherSuite"` // 加密套件
}

type SSLChecker struct {
	Timeout     time.Duration
	WorkerCount int
}

func NewSSLChecker() *SSLChecker {
	return &SSLChecker{
		Timeout:     30 * time.Second, // 增加超时时间
		WorkerCount: runtime.NumCPU(),
	}
}

// 解析密钥用途
func parseKeyUsage(usage x509.KeyUsage) []string {
	var usages []string
	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "数字签名")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "不可否认")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "密钥加密")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "数据加密")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "密钥协商")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "证书签名")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRL签名")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "仅加密")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "仅解密")
	}
	return usages
}

// 获取证书信息
func getCertInfo(cert *x509.Certificate) CertInfo {
	return CertInfo{
		Subject:        cert.Subject.String(),
		Issuer:         cert.Issuer.String(),
		SerialNumber:   cert.SerialNumber.String(),
		Version:        cert.Version,
		IsCA:           cert.IsCA,
		NotBefore:      cert.NotBefore,
		NotAfter:       cert.NotAfter,
		SignatureAlg:   cert.SignatureAlgorithm.String(),
		PublicKeyAlg:   cert.PublicKeyAlgorithm.String(),
		KeyUsage:       parseKeyUsage(cert.KeyUsage),
		DNSNames:       cert.DNSNames,
		EmailAddresses: cert.EmailAddresses,
	}
}

// 并发检查多个域名
func (c *SSLChecker) CheckMultipleSSL(domains []string) []CertResult {
	resultChan := make(chan CertResult, len(domains))
	var wg sync.WaitGroup

	// 创建工作池
	workChan := make(chan string, len(domains))

	// 启动工作协程
	for i := 0; i < c.WorkerCount; i++ {
		wg.Add(1)
		go c.worker(workChan, resultChan, &wg)
	}

	// 分发任务
	for _, domain := range domains {
		workChan <- domain
	}
	close(workChan)

	// 等待所有工作完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	var results []CertResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// 工作协程
func (c *SSLChecker) worker(domains <-chan string, results chan<- CertResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for domain := range domains {
		results <- c.CheckSSLExpiry(domain)
	}
}

func (c *SSLChecker) CheckSSLExpiry(domain string) CertResult {
	result := CertResult{
		Domain: domain,
	}

	conf := &tls.Config{
		ServerName: domain,
		// 先允许跳过验证，我们将在后续手动验证证书
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	var d net.Dialer
	netConn, err := d.DialContext(ctx, "tcp", domain+":443")
	if err != nil {
		result.Error = fmt.Sprintf("连接失败: %v", err)
		return result
	}
	defer netConn.Close()

	conn := tls.Client(netConn, conf)
	if err := conn.Handshake(); err != nil {
		result.Error = fmt.Sprintf("TLS握手失败: %v", err)
		return result
	}
	defer conn.Close()

	state := conn.ConnectionState()
	certs := state.PeerCertificates
	if len(certs) == 0 {
		result.Error = "未找到证书"
		return result
	}

	// 手动验证证书
	var validCert bool
	for _, cert := range certs {
		// 检查 DNS 名称中的泛域名
		for _, dnsName := range cert.DNSNames {
			// 输出调试信息
			fmt.Printf("检查证书域名: %s 与目标域名: %s\n", dnsName, domain)
			if matchWildcard(dnsName, domain) {
				validCert = true
				break
			}
		}
		// 也检查 CommonName
		if matchWildcard(cert.Subject.CommonName, domain) {
			validCert = true
		}
		if validCert {
			break
		}
	}

	if !validCert {
		result.Error = fmt.Sprintf("未找到匹配域名 %s 的证书", domain)
		return result
	}

	// 使用第一个证书作为主证书
	cert := certs[0]
	result.ExpiryDate = cert.NotAfter
	result.RemainingDays = int(time.Until(cert.NotAfter).Hours() / 24)
	result.Status = time.Now().Before(cert.NotAfter)
	result.Protocol = getProtocolVersion(state.Version)
	result.CipherSuite = tls.CipherSuiteName(state.CipherSuite)

	// 获取证书链
	result.CertChain = make([]CertInfo, len(certs))
	for i, cert := range certs {
		result.CertChain[i] = getCertInfo(cert)
	}

	return result
}

// 获取协议版本
func getProtocolVersion(version uint16) string {
	switch version {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return "Unknown"
	}
}

// 检查泛域名匹配
func matchWildcard(pattern, domain string) bool {
	// 直接匹配
	if pattern == domain {
		return true
	}

	// 提取域名的各个部分
	domainParts := strings.Split(domain, ".")
	patternParts := strings.Split(pattern, ".")

	// 如果是泛域名证书
	if len(patternParts) > 0 && patternParts[0] == "*" {
		// 移除第一个部分（通配符）后比较剩余部分
		baseDomain := strings.Join(patternParts[1:], ".")
		if len(domainParts) > len(patternParts)-1 {
			checkDomain := strings.Join(domainParts[1:], ".")
			return checkDomain == baseDomain
		}
	}

	return false
}
