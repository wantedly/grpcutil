package grpcutil

import (
	"context"

	"google.golang.org/grpc/stats"
)

type chainedStatsHandler struct {
	handlers []stats.Handler
}

// ChainStatsHandlers creates a single stats handler out of a chain of many stats handlers.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
//
// See https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware#ChainUnaryServer
func ChainStatsHandlers(handlers ...stats.Handler) stats.Handler {
	return &chainedStatsHandler{
		handlers: handlers,
	}
}

func (csh *chainedStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	for _, h := range csh.handlers {
		ctx = h.TagRPC(ctx, info)
	}
	return ctx
}

func (csh *chainedStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	for _, h := range csh.handlers {
		h.HandleRPC(ctx, s)
	}
}

func (csh *chainedStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	for _, h := range csh.handlers {
		ctx = h.TagConn(ctx, info)
	}
	return ctx
}

func (csh *chainedStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	for _, h := range csh.handlers {
		h.HandleConn(ctx, s)
	}
}
