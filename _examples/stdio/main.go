package main

import (
	"context"
	"fmt"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server"
)

func to[T any](x T) *T {
	return &x
}

func main() {
	prompt := apis.Prompt{
		Name:        "hello",
		Description: to("hello world"),
		Arguments: []apis.PromptArgument{
			{
				Name:        "name",
				Description: to("your name"),
				Required:    to(false),
			},
		},
	}
	promptFunc := func(msg apis.GetPromptRequest) apis.GetPromptResult {
		return apis.GetPromptResult{
			Description: to("Hello MCP"),
			Messages: []apis.PromptMessage{
				{
					Role: apis.Role("user"),
					Content: map[string]any{
						"type": "text",
						"text": fmt.Sprintf("Please respond with your name: %s", msg.Params.Arguments["name"]),
					},
				},
			},
		}
	}
	tool := apis.Tool{
		Name:        "hello",
		Description: to("hello world"),
		InputSchema: apis.ToolInputSchema{
			Type: "object",
			Properties: map[string]map[string]any{
				"name": {
					"type":        "string",
					"description": "your name",
				},
			},
			Required: []string{"name"},
		},
	}
	toolFunc := func(msg apis.CallToolRequest) apis.CallToolResult {
		return apis.CallToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": fmt.Sprintf("Hello, %s!", msg.Params.Arguments["name"]),
				},
			},
			IsError: to(false),
		}
	}
	srv := server.NewServer("stdio", "1.0.0").
		Prompt(&prompt, promptFunc).
		Tool(&tool, toolFunc).
		Build()
	ctx, stdio := server.NewStdioServer(context.Background(), srv)
	stdio.Start(ctx)
}
