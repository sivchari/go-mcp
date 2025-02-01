package server

import (
	"log/slog"
	"os"
	"slices"

	"github.com/google/uuid"

	"github.com/sivchari/go-mcp/internal/apis"
)

type Server struct {
	name    string
	version string
	logger  *slog.Logger

	capabilities *apis.ServerCapabilities
	instructions *string
	prompts      prompts
	tools        tools
}

type prompts struct {
	lists map[string][]apis.Prompt
	funcs map[string]PromptFunc
}

type tools struct {
	lists map[string][]apis.Tool
	funcs map[string]ToolFunc
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
	Tool(tool *apis.Tool, fn ToolFunc) Builder
	PageSize(pageSize int) Builder
}

var _ Builder = &parameter{}

type parameter struct {
	name    string
	version string

	capabilities *apis.ServerCapabilities
	instructions string
	prompts      []apis.Prompt
	promptFuncs  map[string]PromptFunc
	tools        []apis.Tool
	toolsFuncs   map[string]ToolFunc
	pageSize     int
}

func (p *parameter) Build() *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	if p.pageSize == 0 {
		p.pageSize = 10
	}
	promptLists := make(map[string][]apis.Prompt)
	for _, prompt := range slices.Collect(slices.Chunk(p.prompts, p.pageSize)) {
		uid := uuid.New().String()
		promptLists[uid] = prompt
	}
	prompts := prompts{
		lists: promptLists,
		funcs: p.promptFuncs,
	}
	toolLists := make(map[string][]apis.Tool)
	for _, tool := range slices.Collect(slices.Chunk(p.tools, p.pageSize)) {
		uid := uuid.New().String()
		toolLists[uid] = tool
	}
	tools := tools{
		lists: toolLists,
		funcs: p.toolsFuncs,
	}
	return &Server{
		name:         p.name,
		version:      p.version,
		logger:       logger,
		capabilities: p.capabilities,
		instructions: &p.instructions,
		prompts:      prompts,
		tools:        tools,
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

// PromptFunc is a function that returns a prompt result.
//
// You should return messages that all required arguments are filled.
type PromptFunc func(msg apis.GetPromptRequest) apis.GetPromptResult

func (p *parameter) Prompt(prompt *apis.Prompt, fn PromptFunc) Builder {
	if prompt == nil || fn == nil {
		return p
	}
	p.prompts = append(p.prompts, *prompt)
	if p.promptFuncs == nil {
		p.promptFuncs = make(map[string]PromptFunc)
	}
	p.promptFuncs[prompt.Name] = fn
	return p
}

// ToolFunc is a function that returns a tool result.
//
// You should return messages that all required arguments are filled.
type ToolFunc func(msg apis.CallToolRequest) apis.CallToolResult

func (p *parameter) Tool(tool *apis.Tool, fn ToolFunc) Builder {
	if tool == nil || fn == nil {
		return p
	}
	p.tools = append(p.tools, *tool)
	if p.toolsFuncs == nil {
		p.toolsFuncs = make(map[string]ToolFunc)
	}
	p.toolsFuncs[tool.Name] = fn
	return p
}

func (p *parameter) PageSize(pageSize int) Builder {
	p.pageSize = pageSize
	return p
}
