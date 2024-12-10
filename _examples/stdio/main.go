package main

import (
	"context"

	"github.com/sivchari/go-mcp/server"
)

func main() {
	srv := server.NewServer("stdio", "1.0.0")
	ctx, stdio := server.NewStdioServer(context.Background(), srv)
	stdio.Start(ctx)
}
