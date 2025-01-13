package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Log       LogConfig       `yaml:"log"`
	SSL       SSLConfig       `yaml:"ssl"`
	Scanner   ScannerConfig   `yaml:"scanner"`
	WebSocket WebSocketConfig `yaml:"websocket"`
	DNS       DNSConfig       `yaml:"dns"`
	Proxy     ProxyConfig     `yaml:"proxy"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type LogConfig struct {
	Level         string `yaml:"level"`
	Path          string `yaml:"path"`
	Filename      string `yaml:"filename"`
	ErrorFilename string `yaml:"errorFilename"`
	Console       bool   `yaml:"console"`
	MaxAge        int    `yaml:"maxAge"`
	MaxSize       int    `yaml:"maxSize"`
	Compress      bool   `yaml:"compress"`
}

type SSLConfig struct {
	CheckInterval time.Duration `yaml:"checkInterval"`
	ExpiryWarning time.Duration `yaml:"expiryWarning"`
	Timeout       time.Duration `yaml:"timeout"`
	Concurrent    int           `yaml:"concurrent"`
}

type ScannerConfig struct {
	Timeout      time.Duration `yaml:"timeout"`
	Concurrent   int           `yaml:"concurrent"`
	DefaultPorts []int         `yaml:"defaultPorts"`
}

type WebSocketConfig struct {
	PingInterval time.Duration `yaml:"pingInterval"`
	Timeout      time.Duration `yaml:"timeout"`
	RetryTimes   int           `yaml:"retryTimes"`
}

var GlobalConfig Config

// 解析时间字符串为 Duration
func parseDuration(str string) (time.Duration, error) {
	return time.ParseDuration(str)
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 首先解析到临时结构体
	var tmpConfig struct {
		Server ServerConfig `yaml:"server"`
		Log    LogConfig    `yaml:"log"`
		SSL    struct {
			CheckInterval string `yaml:"checkInterval"`
			ExpiryWarning string `yaml:"expiryWarning"`
			Timeout       string `yaml:"timeout"`
			Concurrent    int    `yaml:"concurrent"`
		} `yaml:"ssl"`
		Scanner struct {
			Timeout      string `yaml:"timeout"`
			Concurrent   int    `yaml:"concurrent"`
			DefaultPorts []int  `yaml:"defaultPorts"`
		} `yaml:"scanner"`
		WebSocket struct {
			PingInterval string `yaml:"pingInterval"`
			Timeout      string `yaml:"timeout"`
			RetryTimes   int    `yaml:"retryTimes"`
		} `yaml:"websocket"`
	}

	if err := yaml.Unmarshal(data, &tmpConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 转换时间字符串为 Duration
	checkInterval, err := parseDuration(tmpConfig.SSL.CheckInterval)
	if err != nil {
		return fmt.Errorf("解析 SSL checkInterval 失败: %v", err)
	}

	expiryWarning, err := parseDuration(tmpConfig.SSL.ExpiryWarning)
	if err != nil {
		return fmt.Errorf("解析 SSL expiryWarning 失败: %v", err)
	}

	sslTimeout, err := parseDuration(tmpConfig.SSL.Timeout)
	if err != nil {
		return fmt.Errorf("解析 SSL timeout 失败: %v", err)
	}

	scannerTimeout, err := parseDuration(tmpConfig.Scanner.Timeout)
	if err != nil {
		return fmt.Errorf("解析 Scanner timeout 失败: %v", err)
	}

	wsPingInterval, err := parseDuration(tmpConfig.WebSocket.PingInterval)
	if err != nil {
		return fmt.Errorf("解析 WebSocket pingInterval 失败: %v", err)
	}

	wsTimeout, err := parseDuration(tmpConfig.WebSocket.Timeout)
	if err != nil {
		return fmt.Errorf("解析 WebSocket timeout 失败: %v", err)
	}

	// 设置全局配置
	GlobalConfig = Config{
		Server: tmpConfig.Server,
		Log:    tmpConfig.Log,
		SSL: SSLConfig{
			CheckInterval: checkInterval,
			ExpiryWarning: expiryWarning,
			Timeout:       sslTimeout,
			Concurrent:    tmpConfig.SSL.Concurrent,
		},
		Scanner: ScannerConfig{
			Timeout:      scannerTimeout,
			Concurrent:   tmpConfig.Scanner.Concurrent,
			DefaultPorts: tmpConfig.Scanner.DefaultPorts,
		},
		WebSocket: WebSocketConfig{
			PingInterval: wsPingInterval,
			Timeout:      wsTimeout,
			RetryTimes:   tmpConfig.WebSocket.RetryTimes,
		},
	}

	return nil
}

type ProxyConfig struct {
	Enable   bool          `yaml:"enable"`
	Type     string        `yaml:"type"` // 添加代理类型
	URL      string        `yaml:"url"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

type DNSConfig struct {
	Servers []string      `yaml:"servers"`
	Timeout time.Duration `yaml:"timeout"`
}
