package server

import (
	"maps"
	"slices"

	"github.com/sivchari/go-mcp/apis"
)

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

func (s *Server) cursorTools(cursor *string) []apis.Tool {
	if len(s.tools.lists) == 0 {
		return nil
	}
	if cursor == nil || *cursor == "" {
		return s.tools.lists[slices.Sorted(maps.Keys(s.tools.lists))[0]]
	}
	return s.tools.lists[*cursor]
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
