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
				Required:    to(true),
			},
		},
	}
	promptFunc := func(msg apis.GetPromptRequest) apis.GetPromptResult {
		return apis.GetPromptResult{
			Description: to("Hello MCP"),
			Messages: []apis.PromptMessage{
				{
					Role: apis.Role("user"),
					Content: map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf("Please respond with your name: %s", msg.Params.Arguments["name"]),
					},
				},
			},
		}
	}
	srv := server.NewServer("stdio", "1.0.0").Prompt(&prompt, promptFunc).Build()
	ctx, stdio := server.NewStdioServer(context.Background(), srv)
	stdio.Start(ctx)
}
