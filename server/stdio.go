package server

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
)

func NewStdioServer(ctx context.Context, server *Server) (context.Context, *StdioServer) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	return ctx, &StdioServer{
		Server:     server,
		cancelFunc: cancel,
	}
}

type StdioServer struct {
	Server     *Server
	cancelFunc context.CancelFunc
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
				s.Server.logger.Error("failed to unmarshal request", slog.String("err", err.Error()))
			}
			resp := s.Server.Handle(ctx, msg)
			if resp == nil {
				continue
			}
			if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
				s.Server.logger.Error("failed to marshal response", slog.String("err", err.Error()))
			}
		}
	}
}
