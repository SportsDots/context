package sportctx

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
)

// Context is a structure which wraps context and provide additional methods.
type Context struct {
	context.Context
	requestID     string
	throughParams map[string]string
	logger        *zerolog.Logger
}

type Opt func(c *Context)

// WithContext is an option for adding context.
func WithContext(ctx context.Context) func(c *Context) {
	return func(c *Context) {
		c.Context = ctx
	}
}

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if sportctx, ok := ctx.(*Context); ok {
		var (
			cancelFunc context.CancelFunc
		)

		sportctx.Context, cancelFunc = context.WithTimeout(sportctx.Context, timeout)
		return sportctx, cancelFunc
	}

	return context.WithTimeout(ctx, timeout)
}

func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	if sportctx, ok := ctx.(*Context); ok {
		var (
			cancelFunc context.CancelFunc
		)

		sportctx.Context, cancelFunc = context.WithCancel(sportctx.Context)
		return sportctx, cancelFunc
	}

	return context.WithCancel(ctx)
}

func WithCancelWithoutOverriding(ctx context.Context) (context.Context, context.CancelFunc) {
	if sportctx, ok := ctx.(*Context); ok {
		var (
			cancelFunc context.CancelFunc
		)

		cc := *sportctx
		cc.Context, cancelFunc = context.WithCancel(sportctx.Context)
		return &cc, cancelFunc
	}

	return context.WithCancel(ctx)
}

// WithValue returns a copy of parent in which the value associated with key is
// val.
// The provided key must be comparable and should not be of type
// string or any other built-in type to avoid collisions between
// packages using context. context keys often have concrete type
// struct{}.
func WithValue(ctx context.Context, key, val any) context.Context {
	if sportctx, ok := ctx.(*Context); ok {
		sportctx.Context = context.WithValue(sportctx.Context, key, val)
		return sportctx
	}

	return context.WithValue(ctx, key, val)
}

// WithRequestID is an option for adding specific requestID.
func WithRequestID(requestID string) func(c *Context) {
	return func(c *Context) {
		if requestID != "" {
			c.requestID = requestID
		}
	}
}

// WithThroughParams is an option for adding through params for requests.
func WithThroughParams(values map[string]string) func(c *Context) {
	return func(c *Context) {
		c.throughParams = values
	}
}

// WithLogger is an option for adding logger into context.
func WithLogger(logger *zerolog.Logger) func(c *Context) {
	return func(c *Context) {
		c.logger = logger
	}
}

const (
	xRequestID = "X-Request-ID"
)

// NewContext is a constructor for Context.
func NewContext(opts ...Opt) *Context {
	ctx := getDefaultContext()

	for _, opt := range opts {
		opt(ctx)
	}

	/*
		if we don't have a requestID parameter on input request, then we need generate it in that library and
		put it into throughParams for future relations.
	*/
	if ctx.throughParams == nil {
		ctx.throughParams = make(map[string]string)
	}
	ctx.throughParams[strings.ToLower(xRequestID)] = ctx.requestID

	if ctx.logger != nil {
		log := ctx.logger.
			With().
			Str("request_id", ctx.requestID).
			Logger()
		ctx.logger = &log
	}

	for k, v := range ctx.throughParams {
		ctx.Context = metadata.AppendToOutgoingContext(ctx.Context, k, v)
	}

	return ctx
}

// NewContextWithTimeout is a constructor for Context with timeout.
func NewContextWithTimeout(timeout time.Duration, opts ...Opt) (*Context, context.CancelFunc) {
	ctx := NewContext(opts...)

	var (
		cancel context.CancelFunc
	)

	ctx.Context, cancel = context.WithTimeout(ctx.Context, timeout)
	return ctx, cancel
}

// GetRequestID is a method that return requestID.
func (c *Context) GetRequestID() string {
	return c.requestID
}

// GetThroughParams is a method that return through params.
func (c *Context) GetThroughParams() map[string]string {
	return c.throughParams
}

// GetLogger is a method for getting logger.
func (c *Context) GetLogger() *zerolog.Logger {
	return c.logger
}

func getDefaultContext() *Context {
	ctx := &Context{}

	ctx.Context = context.Background()

	return ctx
}
