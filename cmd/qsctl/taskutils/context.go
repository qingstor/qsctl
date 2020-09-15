package taskutils

import "context"

// ctxKey used to set value in context
var ctxKey struct{}

// ContextWithHandler return a ctx contains handler
func ContextWithHandler(ctx context.Context, handler *PbHandler) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if handler == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKey, handler)
}

// HandlerFromContext get handler from ctx, if not exist, return nil
func HandlerFromContext(ctx context.Context) *PbHandler {
	h, ok := ctx.Value(ctxKey).(*PbHandler)
	if !ok {
		return nil
	}
	return h
}
