package handler

import (
	"context"
	"isscan/internal/pkg/errors"
	"isscan/internal/service"

	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type Handler struct {
	pingService service.PingService
	portService service.PortService
	sslService  service.SSLService
	wsService   service.WebSocketService
}

func NewHandler(factory service.ServiceFactory) *Handler {
	return &Handler{
		pingService: factory.NewPingService(),
		portService: factory.NewPortService(),
		sslService:  factory.NewSSLService(),
		wsService:   factory.NewWebSocketService(),
	}
}

// PingHandler handles ping requests
func (h *Handler) PingHandler(ctx context.Context, c *app.RequestContext) {
	var req PingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, errors.NewError(errors.ErrBadRequest, "无效的请求参数"))
		return
	}

	result, err := h.pingService.Ping(ctx, req.Target)
	if err != nil {
		c.JSON(500, errors.NewError(errors.ErrInternal, err.Error()))
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    result,
	})
}

// ScanPortHandler 处理端口扫描请求
func (h *Handler) ScanPortHandler(ctx context.Context, c *app.RequestContext) {
	var req PortScanRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, errors.NewError(errors.ErrBadRequest, err.Error()))
		return
	}

	// 设置扫描超时
	scanCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	results, err := h.portService.ScanPorts(scanCtx, req.Host, req.Ports)
	if err != nil {
		c.JSON(500, errors.NewError(errors.ErrInternal, err.Error()))
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    results,
	})
}

// CheckSSLHandler 处理SSL证书检查请求
func (h *Handler) CheckSSLHandler(ctx context.Context, c *app.RequestContext) {
	var req SSLCheckRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, errors.NewError(errors.ErrBadRequest, "无效的请求参数"))
		return
	}

	result, err := h.sslService.CheckDomain(ctx, req.Domain)
	if err != nil {
		c.JSON(500, errors.NewError(errors.ErrInternal, err.Error()))
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    result,
	})
}

// CheckWSHandler 处理WebSocket检查请求
func (h *Handler) CheckWSHandler(ctx context.Context, c *app.RequestContext) {
	var req WebSocketCheckRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, errors.NewError(errors.ErrBadRequest, "无效的请求参数"))
		return
	}

	result, err := h.wsService.CheckEndpoint(ctx, req.Host, req.Path, req.Secure)
	if err != nil {
		c.JSON(500, errors.NewError(errors.ErrInternal, err.Error()))
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    result,
	})
}

// BatchCheckSSLHandler 处理批量SSL证书检查请求
func (h *Handler) BatchCheckSSLHandler(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Domains []string `json:"domains"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, errors.NewError(errors.ErrBadRequest, "无效的请求参数"))
		return
	}

	results, err := h.sslService.BatchCheck(ctx, req.Domains)
	if err != nil {
		c.JSON(500, errors.NewError(errors.ErrInternal, err.Error()))
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    results,
	})
}

// BatchCheckWSHandler 处理批量WebSocket检查请求
func (h *Handler) BatchCheckWSHandler(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Endpoints []service.Endpoint `json:"endpoints"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, errors.NewError(errors.ErrBadRequest, "无效的请求参数"))
		return
	}

	results, err := h.wsService.BatchCheck(ctx, req.Endpoints)
	if err != nil {
		c.JSON(500, errors.NewError(errors.ErrInternal, err.Error()))
		return
	}

	c.JSON(200, Response{
		Success: true,
		Data:    results,
	})
}

// ... other handler methods ...
