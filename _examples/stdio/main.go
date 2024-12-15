package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sivchari/go-mcp"
	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server"
)

func main() {
	prompt := apis.Prompt{
		Name:        "hello",
		Description: mcp.Ptr("hello world"),
		Arguments: []apis.PromptArgument{
			{
				Name:        "name",
				Description: mcp.Ptr("your name"),
				Required:    mcp.Ptr(false),
			},
		},
	}
	promptFunc := func(msg apis.GetPromptRequest) apis.GetPromptResult {
		return apis.GetPromptResult{
			Description: mcp.Ptr("Hello MCP"),
			Messages: []apis.PromptMessage{
				{
					Role:    apis.RoleUser,
					Content: mcp.NewTextContent("Please respond with your name"),
				},
				{
					Role: apis.RoleAssistant,
					Content: mcp.NewResource(apis.TextResourceContents{
						Text:     "Hello, what is your name?",
						Uri:      "file://your-name.txt",
						MimeType: mcp.Ptr("text/plain"),
					}),
				},
			},
		}
	}

	tool := apis.Tool{
		Name:        "hello",
		Description: mcp.Ptr("hello world"),
		InputSchema: *mcp.NewToolInput().
			WithString("name").
			WithNumber("age").
			WithRequired("name", "age").
			Build(),
	}
	toolFunc := func(msg apis.CallToolRequest) apis.CallToolResult {
		return apis.CallToolResult{
			Content: []any{
				mcp.NewTextContent(fmt.Sprintf("Hello %s, you are %d years old", msg.Params.Arguments["name"], msg.Params.Arguments["age"])),
			},
			IsError: mcp.Ptr(false),
		}
	}

	srv := server.NewServer("stdio", "1.0.0").
		Prompt(&prompt, promptFunc).
		Tool(&tool, toolFunc).
		Build()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	stdio := server.NewStdioServer(srv)
	stdio.Start(ctx)
	<-ctx.Done()
}
