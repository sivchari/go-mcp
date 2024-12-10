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
		s.logger.Error("failed to unmarshal request", slog.String("err", err.Error()))
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	nextCursor := s.nextCursor(req.Params.Cursor)
	prompts := s.cursorPrompts(nextCursor)
	if prompts == nil {
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
		for i, token := range slices.Collect(maps.Keys(s.prompts)) {
			if token == *cursor {
				if i+1 < len(s.prompts) {
					nextCursor = &slices.Collect(maps.Keys(s.prompts))[i+1]
				}
			}
		}
	}
	if nextCursor == nil && len(s.prompts) > 0 {
		nextCursor = &slices.Collect(maps.Keys(s.prompts))[0]
	}
	return nextCursor
}

func (s *Server) cursorPrompts(cursor *string) []apis.Prompt {
	if cursor == nil && len(s.prompts) == 0 {
		return nil
	}
	if cursor == nil {
		return slices.Collect(maps.Values(s.prompts))[0]
	}
	return s.prompts[*cursor]
}
