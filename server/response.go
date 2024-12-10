package server

import (
	"context"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal"
	"github.com/sivchari/go-mcp/server/internal/xcontext"
)

func Error(ctx context.Context, code int, err error) *apis.JSONRPCError {
	if err == nil {
		return nil
	}
	return &apis.JSONRPCError{
		Id:      apis.RequestId(xcontext.ID(ctx)),
		Jsonrpc: internal.JSONRPCVersion,
		Error: apis.JSONRPCErrorError{
			Code: code,
			Data: map[string]interface{}{
				"method": xcontext.Method(ctx),
			},
			Message: err.Error(),
		},
	}
}
