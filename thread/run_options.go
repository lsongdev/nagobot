package thread

import "context"

type sinkControlKey struct{}

// WithoutSink disables sink delivery for this Run call.
func WithoutSink(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, sinkControlKey{}, true)
}

// WithSink forces sink delivery for this Run call.
func WithSink(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, sinkControlKey{}, false)
}

func isSinkSuppressed(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	v, _ := ctx.Value(sinkControlKey{}).(bool)
	return v
}
