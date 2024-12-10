package main

import (
	"context"

	"github.com/sivchari/go-mcp/apis"
	"github.com/sivchari/go-mcp/server"
)

func to[T any](x T) *T {
	return &x
}

func main() {
	prompt := apis.Prompt{
		Name:        "sample",
		Description: to("sample description"),
		Arguments: []apis.PromptArgument{
			{
				Name:        "arg1",
				Description: to("arg1 description"),
				Required:    to(true),
			},
		},
	}
	srv := server.NewServer("stdio", "1.0.0").Prompt(&prompt).Build()
	ctx, stdio := server.NewStdioServer(context.Background(), srv)
	stdio.Start(ctx)
}
