package xcontext

import (
	"context"
)

type idKey struct{}

var idKeyInstance = idKey{}

func WithID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, idKeyInstance, id)
}

func ID(ctx context.Context) int {
	id, ok := ctx.Value(idKeyInstance).(int)
	if !ok {
		return 0
	}
	return id
}

type methodKey struct{}

var methodKeyInstance = methodKey{}

func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKeyInstance, method)
}

func Method(ctx context.Context) string {
	method, ok := ctx.Value(methodKeyInstance).(string)
	if !ok {
		return ""
	}
	return method
}
