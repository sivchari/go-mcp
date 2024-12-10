package server

import (
	"context"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal/xcontext"
)

func Error(ctx context.Context, code int, err error) *apis.JSONRPCErrorError {
	if err == nil || code == 0 {
		return nil
	}
	return &apis.JSONRPCErrorError{
		Code: code,
		Data: map[string]interface{}{
			"method": xcontext.Method(ctx),
		},
		Message: err.Error(),
	}
}
