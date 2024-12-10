package server

import (
	"log/slog"
	"os"
	"slices"

	"github.com/google/uuid"

	"github.com/sivchari/go-mcp/apis"
)

type Server struct {
	name    string
	version string
	logger  *slog.Logger

	capabilities *apis.ServerCapabilities
	instructions *string
	prompts      map[string][]apis.Prompt
}

func NewServer(name, version string) Builder {
	return &parameter{
		name:    name,
		version: version,
	}
}

type Builder interface {
	Build() *Server
	Capabilities(capabilities *apis.ServerCapabilities) Builder
	Instructions(instructions string) Builder
	Prompt(prompt *apis.Prompt) Builder
	PageSize(pageSize int) Builder
}

var _ Builder = &parameter{}

type parameter struct {
	name    string
	version string

	capabilities *apis.ServerCapabilities
	instructions string
	prompts      []apis.Prompt
	pageSize     int
}

func (p *parameter) Build() *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if p.pageSize == 0 {
		p.pageSize = 10
	}
	prompts := make(map[string][]apis.Prompt)
	for _, prompt := range slices.Collect(slices.Chunk(p.prompts, p.pageSize)) {
		uid := uuid.New().String()
		prompts[uid] = prompt
	}
	return &Server{
		name:         p.name,
		version:      p.version,
		logger:       logger,
		capabilities: p.capabilities,
		instructions: &p.instructions,
		prompts:      prompts,
	}
}

func (p *parameter) Capabilities(capabilities *apis.ServerCapabilities) Builder {
	p.capabilities = capabilities
	return p
}

func (p *parameter) Instructions(instructions string) Builder {
	p.instructions = instructions
	return p
}

func (p *parameter) Prompt(prompt *apis.Prompt) Builder {
	if prompt == nil {
		return p
	}
	p.prompts = append(p.prompts, *prompt)
	return p
}

func (p *parameter) PageSize(pageSize int) Builder {
	p.pageSize = pageSize
	return p
}
