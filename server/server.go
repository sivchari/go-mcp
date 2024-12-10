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
	prompts      prompts
}

type prompts struct {
	lists map[string][]apis.Prompt
	funcs map[string]PromptFunc
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
	Prompt(prompt *apis.Prompt, fn PromptFunc) Builder
	PageSize(pageSize int) Builder
}

var _ Builder = &parameter{}

type parameter struct {
	name    string
	version string

	capabilities *apis.ServerCapabilities
	instructions string
	prompts      []apis.Prompt
	funcs        map[string]PromptFunc
	pageSize     int
}

func (p *parameter) Build() *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if p.pageSize == 0 {
		p.pageSize = 10
	}
	lists := make(map[string][]apis.Prompt)
	for _, prompt := range slices.Collect(slices.Chunk(p.prompts, p.pageSize)) {
		uid := uuid.New().String()
		lists[uid] = prompt
	}
	prompts := prompts{
		lists: lists,
		funcs: p.funcs,
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

type PromptFunc func(msg apis.GetPromptRequest) apis.GetPromptResult

func (p *parameter) Prompt(prompt *apis.Prompt, fn PromptFunc) Builder {
	if prompt == nil {
		return p
	}
	p.prompts = append(p.prompts, *prompt)
	if p.funcs == nil {
		p.funcs = make(map[string]PromptFunc)
	}
	p.funcs[prompt.Name] = fn
	return p
}

func (p *parameter) PageSize(pageSize int) Builder {
	p.pageSize = pageSize
	return p
}
