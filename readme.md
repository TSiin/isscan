# ISScan - 基础设施扫描工具 🚀

[![Go Version](https://img.shields.io/badge/Go-1.23.4-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> 一个强大的基础设施扫描工具，提供 SSL 证书检查、端口扫描、Ping 测试和 WebSocket 检测等功能。

## ✨ 功能特性

- 🔒 SSL 证书检查
  - 支持泛域名证书验证
  - 证书过期预警
  - 证书链验证
- 🔍 端口扫描
  - 自定义端口范围
  - 服务识别
  - 并发扫描
- 📡 网络检测
  - Ping 延迟测试
  - TCP 连接测试
  - WebSocket 可用性检测
- 🌐 DNS 解析
  - 支持 DoH (DNS over HTTPS)
  - 多 DNS 服务器支持
  - IPv4/IPv6 解析
- 📊 监控集成
  - Prometheus 指标导出
  - Grafana 仪表盘
  - 健康检查接口

## 🚀 快速开始

### Docker 部署

```bash
# 克隆仓库
git clone https://github.com/TSiin/isscan.git
cd isscan

# 使用 Docker Compose 启动服务
docker-compose up -d
```

### 手动部署

```bash
# 编译
go build -o isscan

# 运行
./isscan --config config.yaml
```

## 📡 API 接口

### 1. 端口扫描

检查指定主机的端口开放情况。

**请求:**
```bash
POST /api/port/scan
Content-Type: application/json

{
    "target": "example.com",
    "ports": [80, 443, 22]
}
```

**响应:**
```json
{
    "success": true,
    "data": [
        {
            "port": 80,
            "open": true,
            "service": "HTTP",
            "latency": "45.2ms"
        },
        {
            "port": 443,
            "open": true,
            "service": "HTTPS",
            "latency": "46.8ms"
        }
    ]
}
```

### 2. SSL 证书检查

检查指定域名的 SSL 证书信息。

**请求:**
```bash
POST /api/ssl/check
Content-Type: application/json

{
    "domain": "example.com",
    "port": 443
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "domain": "example.com",
        "valid": true,
        "issuer": "Let's Encrypt Authority X3",
        "subject": "example.com",
        "notBefore": "2024-01-01T00:00:00Z",
        "notAfter": "2024-03-31T23:59:59Z",
        "dnsNames": ["example.com", "*.example.com"],
        "serialNumber": "1234567890",
        "version": 3,
        "signatureAlgorithm": "SHA256-RSA"
    }
}
```

### 3. Ping 测试

对指定主机执行 ICMP ping 测试。

**请求:**
```bash
POST /api/ping
Content-Type: application/json

{
    "target": "example.com",
    "count": 4,
    "timeout": "2s"
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "host": "example.com",
        "ip": "93.184.216.34",
        "sent": 4,
        "received": 4,
        "loss": 0,
        "minLatency": "45.2ms",
        "avgLatency": "46.8ms",
        "maxLatency": "48.1ms"
    }
}
```

### 4. WebSocket 检测

检查 WebSocket 服务的可用性。

**请求:**
```bash
POST /api/websocket/check
Content-Type: application/json

{
    "url": "ws://example.com/ws",
    "protocols": ["v1"],
    "headers": {
        "Authorization": "Bearer token"
    }
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "connected": true,
        "latency": "35.6ms",
        "protocols": ["v1"],
        "extensions": []
    }
}
```

### 5. DNS 解析

使用 DoH (DNS over HTTPS) 进行域名解析。

**请求:**
```bash
POST /api/dns/resolve
Content-Type: application/json

{
    "domain": "example.com",
    "type": "A",
    "server": "https://cloudflare-dns.com/dns-query"
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "domain": "example.com",
        "records": [
            {
                "type": "A",
                "value": "93.184.216.34",
                "ttl": 300
            }
        ],
        "server": "cloudflare-dns.com",
        "latency": "86.4ms"
    }
}
```

## 🛠 高级功能

### 批量扫描

支持批量扫描多个目标。

**请求:**
```bash
POST /api/batch/scan
Content-Type: application/json

{
    "targets": [
        {
            "host": "example1.com",
            "ports": [80, 443]
        },
        {
            "host": "example2.com",
            "ports": [22, 3306]
        }
    ],
    "concurrent": 2,
    "timeout": "30s"
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "total": 2,
        "completed": 2,
        "results": [
            {
                "host": "example1.com",
                "scans": [
                    {
                        "port": 80,
                        "open": true,
                        "service": "HTTP",
                        "latency": "45.2ms"
                    }
                ]
            }
        ]
    }
}
```

### 证书链验证

完整的 SSL 证书链验证。

**请求:**
```bash
POST /api/ssl/verify-chain
Content-Type: application/json

{
    "domain": "example.com",
    "port": 443,
    "validateChain": true
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "valid": true,
        "chain": [
            {
                "subject": "example.com",
                "issuer": "Let's Encrypt Authority X3",
                "validFrom": "2024-01-01T00:00:00Z",
                "validTo": "2024-03-31T23:59:59Z"
            },
            {
                "subject": "Let's Encrypt Authority X3",
                "issuer": "DST Root CA X3",
                "validFrom": "2016-03-17T16:40:46Z",
                "validTo": "2021-03-17T16:40:46Z"
            }
        ],
        "verificationDetails": {
            "nameConstraints": true,
            "keyUsage": true,
            "extendedKeyUsage": true,
            "basicConstraints": true
        }
    }
}
```

### TCP 连接测试

测试 TCP 连接的建立时间和可靠性。

**请求:**
```bash
POST /api/tcp/test
Content-Type: application/json

{
    "host": "example.com",
    "port": 80,
    "count": 3,
    "interval": "1s",
    "timeout": "5s"
}
```

**响应:**
```json
{
    "success": true,
    "data": {
        "successful": 3,
        "failed": 0,
        "avgConnectTime": "35.6ms",
        "minConnectTime": "32.1ms",
        "maxConnectTime": "38.2ms",
        "details": [
            {
                "timestamp": "2024-01-20T10:00:00Z",
                "successful": true,
                "connectTime": "35.6ms"
            }
        ]
    }
}
```

## 🔧 环境要求

- Go 1.23.4 或更高版本
- Docker 20.10.0 或更高版本（如果使用 Docker 部署）
- 操作系统：Linux, macOS, Windows

## 📦 安装说明

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/TSiin/isscan.git

# 进入项目目录
cd isscan

# 安装依赖
go mod download

# 编译
make build
```

## 📊 监控指标

服务暴露了以下 Prometheus 指标：

- `isscan_requests_total{path="/api/*"}` - API 请求总数
- `isscan_request_duration_seconds` - 请求处理时间
- `isscan_scan_total{type="port|ssl|ping|websocket"}` - 扫描次数
- `isscan_scan_errors_total{type="port|ssl|ping|websocket"}` - 扫描错误次数
- `isscan_up` - 服务存活状态

访问 `/metrics` 端点获取完整的指标数据。

## 📝 配置说明

详细的配置项说明请参考 [config.yaml](config.yaml) 文件。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

## 📄 开源协议

本项目采用 MIT 协议开源，详见 [LICENSE](LICENSE) 文件。

