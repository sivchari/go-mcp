package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal"
)

func (s *Server) Initialize(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.InitializeRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request",
			slog.String("method", internal.MethodInitialize),
			slog.String("err", err.Error()),
		)
		return Error(ctx, internal.CodeInvalidParams, err)
	}

	// TODO: Handle meta field.
	result := &apis.InitializeResult{
		ProtocolVersion: req.Params.ProtocolVersion,
		ServerInfo: apis.Implementation{
			Name:    s.name,
			Version: s.version,
		},
	}

	if s.capabilities != nil {
		result.Capabilities = *s.capabilities
	}

	if s.instructions != nil {
		result.Instructions = s.instructions
	}

	return result
}
