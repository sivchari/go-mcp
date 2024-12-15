package server

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/sivchari/go-mcp/server/internal"
	"github.com/sivchari/go-mcp/server/internal/xcontext"
)

func NewStdioServer(server *Server) *StdioServer {
	return &StdioServer{
		Server: server,
	}
}

type StdioServer struct {
	Server *Server
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func (s *StdioServer) Start(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				return
			}
			txt := scanner.Text()
			var msg json.RawMessage
			if err := json.Unmarshal([]byte(txt), &msg); err != nil {
				s.Server.logger.Error("failed to unmarshal request",
					slog.String("err", err.Error()),
				)
			}
			ctx, result := s.Server.Handle(ctx, msg)
			if ctx == nil || result == nil {
				continue
			}
			resp := Response{
				JSONRPC: internal.JSONRPCVersion,
				ID:      xcontext.ID(ctx),
				Result:  result,
			}
			if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
				s.Server.logger.Error("failed to marshal response",
					slog.String("err", err.Error()),
				)
			}
		}
	}
}
