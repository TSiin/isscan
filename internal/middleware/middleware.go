package middleware

import (
	"context"
	"isscan/pkg/logger"
	"time"

	"strconv"

	"isscan/pkg/metrics"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", zap.Any("error", err))
				c.JSON(500, map[string]interface{}{
					"success": false,
					"error":   "Internal server error",
				})
			}
		}()
		c.Next(ctx)
	}
}

func RequestLogger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())

		c.Next(ctx)

		latency := time.Since(start)
		logger.Info("Request completed",
			zap.String("path", path),
			zap.Int("status", c.Response.StatusCode()),
			zap.Duration("latency", latency),
		)
	}
}

func RequestTracing() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		traceID := uuid.New().String()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		c.Next(ctx)
	}
}

func Metrics() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())

		c.Next(ctx)

		duration := time.Since(start)
		status := c.Response.StatusCode()

		metrics.RequestDuration.WithLabelValues(path).Observe(duration.Seconds())
		metrics.RequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
	}
}
