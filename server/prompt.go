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

func (s *Server) nextCursor(cursor *string) *string {
	var nextCursor *string
	if cursor != nil && *cursor != "" {
		for i, token := range slices.Sorted(maps.Keys(s.prompts)) {
			if token == *cursor {
				if i+1 < len(s.prompts) {
					nextCursor = &slices.Sorted(maps.Keys(s.prompts))[i+1]
				}
			}
		}
	}
	if (cursor == nil || *cursor == "") && len(s.prompts) > 1 {
		nextCursor = &slices.Sorted((maps.Keys(s.prompts)))[1]
	}
	return nextCursor
}

func (s *Server) cursorPrompts(cursor *string) []apis.Prompt {
	if len(s.prompts) == 0 {
		return nil
	}
	if cursor == nil || *cursor == "" {
		return slices.Collect(maps.Values(s.prompts))[0]
	}
	return s.prompts[*cursor]
}
