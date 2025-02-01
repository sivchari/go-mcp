package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sivchari/go-mcp/internal/apis"
	"github.com/sivchari/go-mcp/server/internal"
)

func (s *Server) Tools(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.ListToolsRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request",
			slog.String("method", internal.MethodListTools),
			slog.String("err", err.Error()),
		)
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	nextCursor := s.nextCursor(req.Params.Cursor)
	tools := s.cursorTools(req.Params.Cursor)
	if tools == nil {
		if req.Params.Cursor != nil {
			s.logger.Error("invalid cursor",
				slog.String("cursor", *req.Params.Cursor),
				slog.String("method", internal.MethodListTools),
			)
			return Error(ctx, internal.CodeInvalidParams, nil)
		}
		return &apis.ListToolsResult{}
	}
	return &apis.ListToolsResult{
		Tools:      tools,
		NextCursor: nextCursor,
	}
}

func (s *Server) Call(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.CallToolRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request",
			slog.String("method", internal.MethodCallTool),
			slog.String("err", err.Error()),
		)
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	fn, ok := s.tools.funcs[req.Params.Name]
	if !ok {
		s.logger.Error("tool not found",
			slog.String("name", req.Params.Name),
			slog.String("method", internal.MethodCallTool),
		)
		return Error(ctx, internal.CodeInvalidParams, nil)
	}
	response := fn(req)
	return &response
}
