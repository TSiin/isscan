package errors

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

type ErrorCode int

const (
	ErrBadRequest ErrorCode = iota + 400
	ErrUnauthorized
	ErrForbidden
	ErrNotFound
	ErrInternal = 500
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Hertz 错误处理中间件
func ErrorHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		// 检查是否有错误
		if err := c.Errors.Last(); err != nil {
			var apiErr *Error
			if errors.As(err.Err, &apiErr) {
				c.JSON(int(apiErr.Code), Response{
					Success: false,
					Error:   apiErr.Message,
				})
				return
			}

			// 默认内部错误
			c.JSON(int(ErrInternal), Response{
				Success: false,
				Error:   "Internal server error",
			})
		}
	}
}
