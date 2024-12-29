package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal"
)

func (s *Server) Prompts(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.ListPromptsRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request",
			slog.String("method", internal.MethodListPrompts),
			slog.String("err", err.Error()),
		)
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	nextCursor := s.nextCursor(req.Params.Cursor)
	prompts := s.cursorPrompts(req.Params.Cursor)
	if prompts == nil {
		if req.Params.Cursor != nil {
			s.logger.Error("invalid cursor",
				slog.String("cursor", *req.Params.Cursor),
				slog.String("method", internal.MethodListPrompts),
			)
			return Error(ctx, internal.CodeInvalidParams, nil)
		}
		return &apis.ListPromptsResult{}
	}
	return &apis.ListPromptsResult{
		Prompts:    prompts,
		NextCursor: nextCursor,
	}
}

func (s *Server) Prompt(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	// TODO: Missing required arguments: -32602 (Invalid params)
	var req apis.GetPromptRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request",
			slog.String("method", internal.MethodGetPrompt),
			slog.String("err", err.Error()),
		)
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	fn, ok := s.prompts.funcs[req.Params.Name]
	if !ok {
		s.logger.Error("prompt not found",
			slog.String("name", req.Params.Name),
			slog.String("method", internal.MethodGetPrompt),
		)
		return Error(ctx, internal.CodeInvalidParams, nil)
	}
	response := fn(req)
	return &response
}
