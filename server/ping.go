package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sivchari/go-mcp/internal/apis"
	"github.com/sivchari/go-mcp/server/internal"
	"github.com/sivchari/go-mcp/server/internal/xcontext"
)

func (s *Server) Ping(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.PingRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request",
			slog.String("method", internal.MethodPing),
			slog.String("err", err.Error()),
		)
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	resp := &apis.JSONRPCResponse{
		Jsonrpc: internal.JSONRPCVersion,
		Id:      apis.RequestId(xcontext.ID(ctx)),
		Result:  apis.Result{},
	}
	return resp
}
