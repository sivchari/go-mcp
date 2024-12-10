package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal"
	"github.com/sivchari/go-mcp/server/internal/xcontext"
)

func (s *Server) Handle(ctx context.Context, msg json.RawMessage) (context.Context, json.Unmarshaler) {
	var req apis.JSONRPCRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		if s.checkInitializedNotification(msg) {
			// Ignore the error if the message is an initialized notification.
			return ctx, nil
		}
		s.logger.Error("failed to unmarshal request", slog.String("err", err.Error()))
		return nil, Error(ctx, internal.CodeInvalidParams, err)
	}
	idctx := xcontext.WithID(ctx, int(req.Id))
	methodctx := xcontext.WithMethod(idctx, req.Method)
	switch req.Method {
	case internal.MethodPing:
		return methodctx, s.Ping(methodctx, msg)
	case internal.MethodInitialize:
		return methodctx, s.Initialize(methodctx, msg)
	case internal.MethodGetPrompt:
		return methodctx, s.Prompt(methodctx, msg)
	case internal.MethodListPrompts:
		return methodctx, s.Prompts(methodctx, msg)
	default:
		// not implemented yet
		s.logger.Error("method not implemented", slog.String("method", req.Method))
	}
	return nil, Error(ctx, internal.CodeInternalError, nil)
}

func (s *Server) checkInitializedNotification(msg json.RawMessage) bool {
	var req apis.InitializedNotification
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error(
			"failed to unmarshal request for checking if req is notifications/initialized",
			slog.String("err", err.Error()),
		)
		return false
	}
	if req.Method == internal.MethodNotificationsInitialized {
		return true
	}
	return false
}
