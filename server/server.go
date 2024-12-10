package server

import (
	"log/slog"
	"os"

	"github.com/sivchari/go-mcp/apis"
)

type Server struct {
	name    string
	version string
	logger  *slog.Logger

	capabilities apis.ServerCapabilities
	instructions string
}

func NewServer(name, version string) *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return &Server{
		name:    name,
		version: version,
		logger:  logger,
	}
}

type ServerOption func(*Server)

func WithCapabilities(capabilities apis.ServerCapabilities) ServerOption {
	return func(s *Server) {
		s.capabilities = capabilities
	}
}
