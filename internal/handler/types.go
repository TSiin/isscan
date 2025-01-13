package handler

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PingRequest Ping请求结构
type PingRequest struct {
	Target string `json:"target"`
}

// PortScanRequest 端口扫描请求
type PortScanRequest struct {
	Host  string `json:"host"`
	Ports []int  `json:"ports,omitempty"`
}

// SSLCheckRequest SSL检查请求
type SSLCheckRequest struct {
	Domain string `json:"domain"`
}

// WebSocketCheckRequest WebSocket检查请求
type WebSocketCheckRequest struct {
	Host   string `json:"host"`
	Path   string `json:"path"`
	Secure bool   `json:"secure"`
}
