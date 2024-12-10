package server

import (
	"maps"
	"slices"
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
