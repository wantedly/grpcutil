package statsutil

import (
	"context"

	"google.golang.org/grpc/stats"
)

type chainedHandler struct {
	handlers []stats.Handler
}

// ChainHandlers creates a single stats handler out of a chain of many stats handlers.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
//
// See https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware#ChainUnaryServer
func ChainHandlers(handlers ...stats.Handler) stats.Handler {
	return &chainedHandler{
		handlers: handlers,
	}
}

func (csh *chainedHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	for _, h := range csh.handlers {
		ctx = h.TagRPC(ctx, info)
	}
	return ctx
}

func (csh *chainedHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	for _, h := range csh.handlers {
		h.HandleRPC(ctx, s)
	}
}

func (csh *chainedHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	for _, h := range csh.handlers {
		ctx = h.TagConn(ctx, info)
	}
	return ctx
}

func (csh *chainedHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	for _, h := range csh.handlers {
		h.HandleConn(ctx, s)
	}
}
