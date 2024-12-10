package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server/internal"
	"github.com/sivchari/go-mcp/server/internal/ptr"
	"github.com/sivchari/go-mcp/server/internal/xcontext"
)

func (s *Server) Handle(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.JSONRPCRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request", slog.String("err", err.Error()))
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	idctx := xcontext.WithID(ctx, int(req.Id))
	methodctx := xcontext.WithMethod(idctx, req.Method)
	switch req.Method {
	case internal.MethodPing:
		return s.handlePing(methodctx, msg)
	case internal.MethodInitialize:
		return s.handleInitialize(methodctx, msg)
	default:
		// not implemented yet
		s.logger.Error("method not implemented", slog.String("method", req.Method))
	}
	return nil
}

func (s *Server) handlePing(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.PingRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request", slog.String("err", err.Error()))
		return Error(ctx, internal.CodeInvalidParams, err)
	}
	resp := &apis.JSONRPCResponse{
		Jsonrpc: internal.JSONRPCVersion,
		Id:      apis.RequestId(xcontext.ID(ctx)),
		Result:  apis.Result{},
	}
	return resp
}

func (s *Server) handleInitialize(ctx context.Context, msg json.RawMessage) json.Unmarshaler {
	var req apis.InitializeRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		s.logger.Error("failed to unmarshal request", slog.String("err", err.Error()))
		return Error(ctx, internal.CodeInvalidParams, err)
	}

	// TODO: Handle meta field.
	result := apis.InitializeResult{
		ProtocolVersion: req.Params.ProtocolVersion,
		Capabilities:    s.capabilities,
		ServerInfo: apis.Implementation{
			Name:    s.name,
			Version: s.version,
		},
		Instructions: ptr.To(s.instructions),
	}

	return &apis.JSONRPCResponse{
		Jsonrpc: internal.JSONRPCVersion,
		Id:      apis.RequestId(xcontext.ID(ctx)),
		Result: apis.Result{
			AdditionalProperties: result,
		},
	}
}
