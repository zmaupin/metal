package worker

import "context"

// Interface is an interface for doing work within a go-routine. When work is
// done within a go-routine it can be tricky to return errors. Furthermore,
// if we need additional values, we can access them from the type implementing
// the Worker interface instead of relying on access to free variables within
// scope access.
type Interface interface {
	// Work is a blocking function. If you encouter an error within your Work
	// method, simply pass it along on the channel.
	Work(ctx context.Context, ch chan error)
}

// Func implements Worker as a func
type Func func(ctx context.Context, ch chan error)

// Work satisfies the Interface
func (f Func) Work(ctx context.Context, ch chan error) {
	f(ctx, ch)
}
