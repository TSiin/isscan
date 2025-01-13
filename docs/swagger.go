package docs

// @title Server Manager API
// @version 1.0
// @description 服务器管理和监控工具 API 文档
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// Ping 相关接口
// @Summary 测试主机连通性
// @Description 对指定主机执行 ping 测试
// @Tags Ping
// @Accept json
// @Produce json
// @Param request body PingRequest true "Ping请求参数"
// @Success 200 {object} Response{data=PingResult} "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /ping [post]
// @Example 请求示例
// {
//   "target": "example.com"
// }
// @Example 响应示例
// {
//   "success": true,
//   "data": {
//     "target": "example.com",
//     "ipv4": "93.184.216.34",
//     "ipv6": "2606:2800:220:1:248:1893:25c8:1946",
//     "min_latency": "20ms",
//     "max_latency": "25ms",
//     "avg_latency": "22.5ms",
//     "packet_loss": 0,
//     "send_count": 3,
//     "recv_count": 3,
//     "protocol": "IPv4"
//   }
// }

// 端口扫描相关接口
// @Summary 扫描主机端口
// @Description 扫描指定主机的端口状态
// @Tags Port
// @Accept json
// @Produce json
// @Param request body PortScanRequest true "端口扫描请求参数"
// @Success 200 {object} Response{data=[]PortResult} "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /port/scan [post]
// @Example 请求示例
// {
//   "host": "example.com",
//   "ports": [80, 443, 22, 3306]
// }
// @Example 响应示例
// {
//   "success": true,
//   "data": [
//     {
//       "port": 80,
//       "open": true,
//       "service": "HTTP"
//     },
//     {
//       "port": 443,
//       "open": true,
//       "service": "HTTPS"
//     },
//     {
//       "port": 22,
//       "open": false,
//       "error": "connection refused"
//     }
//   ]
// }

// SSL证书检查相关接口
// @Summary 检查SSL证书
// @Description 检查域名的SSL证书状态
// @Tags SSL
// @Accept json
// @Produce json
// @Param request body SSLCheckRequest true "SSL检查请求参数"
// @Success 200 {object} Response{data=SSLResult} "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /ssl/check [post]
// @Example 请求示例
// {
//   "domain": "example.com"
// }
// @Example 响应示例
// {
//   "success": true,
//   "data": {
//     "valid": true,
//     "expire_at": "2024-12-31T23:59:59Z",
//     "dns_names": ["example.com", "*.example.com"],
//     "issuer": "DigiCert Inc"
//   }
// }

// WebSocket检测相关接口
// @Summary 检查WebSocket连接
// @Description 检查WebSocket端点的可用性
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param request body WebSocketCheckRequest true "WebSocket检查请求参数"
// @Success 200 {object} Response{data=WebSocketResult} "成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /websocket/check [post]
// @Example 请求示例
// {
//   "host": "example.com",
//   "path": "/ws",
//   "secure": true
// }
// @Example 响应示例
// {
//   "success": true,
//   "data": {
//     "host": "example.com",
//     "path": "/ws",
//     "protocol": "wss",
//     "available": true,
//     "latency": "150ms"
//   }
// }

// 批量SSL检查接口
// @Summary 批量检查SSL证书
// @Description 批量检查多个域名的SSL证书状态
// @Tags SSL
// @Accept json
// @Produce json
// @Param request body BatchSSLCheckRequest true "批量SSL检查请求参数"
// @Success 200 {object} Response{data=map[string]SSLResult} "成功"
// @Router /ssl/batch [post]
// @Example 请求示例
// {
//   "domains": ["example.com", "google.com"]
// }
// @Example 响应示例
// {
//   "success": true,
//   "data": {
//     "example.com": {
//       "valid": true,
//       "expire_at": "2024-12-31T23:59:59Z",
//       "dns_names": ["example.com"],
//       "issuer": "DigiCert Inc"
//     },
//     "google.com": {
//       "valid": true,
//       "expire_at": "2024-06-30T23:59:59Z",
//       "dns_names": ["google.com", "*.google.com"],
//       "issuer": "Google Trust Services"
//     }
//   }
// }

// 批量WebSocket检查接口
// @Summary 批量检查WebSocket连接
// @Description 批量检查多个WebSocket端点的可用性
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param request body BatchWSCheckRequest true "批量WebSocket检查请求参数"
// @Success 200 {object} Response{data=[]WebSocketResult} "成功"
// @Router /websocket/batch [post]
// @Example 请求示例
// {
//   "endpoints": [
//     {
//       "host": "example.com",
//       "path": "/ws1",
//       "secure": true
//     },
//     {
//       "host": "example.org",
//       "path": "/ws2",
//       "secure": false
//     }
//   ]
// }
// @Example 响应示例
// {
//   "success": true,
//   "data": [
//     {
//       "host": "example.com",
//       "path": "/ws1",
//       "protocol": "wss",
//       "available": true,
//       "latency": "150ms"
//     },
//     {
//       "host": "example.org",
//       "path": "/ws2",
//       "protocol": "ws",
//       "available": false,
//       "error": "connection refused"
//     }
//   ]
// }
