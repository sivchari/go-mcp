package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"slices"

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
	result := fn(req)
	return &result
}

func (s *Server) nextCursor(cursor *string) *string {
	var nextCursor *string
	if cursor != nil && *cursor != "" {
		for i, token := range slices.Sorted(maps.Keys(s.prompts.lists)) {
			if token == *cursor {
				if i+1 < len(s.prompts.lists) {
					nextCursor = &slices.Sorted(maps.Keys(s.prompts.lists))[i+1]
				}
			}
		}
	}
	if (cursor == nil || *cursor == "") && len(s.prompts.lists) > 1 {
		nextCursor = &slices.Sorted((maps.Keys(s.prompts.lists)))[1]
	}
	return nextCursor
}

func (s *Server) cursorPrompts(cursor *string) []apis.Prompt {
	if len(s.prompts.lists) == 0 {
		return nil
	}
	if cursor == nil || *cursor == "" {
		return slices.Collect(maps.Values(s.prompts.lists))[0]
	}
	return s.prompts.lists[*cursor]
}
