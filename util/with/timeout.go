package with

import (
	"context"
	"errors"
	"time"

	"github.com/metal-go/metal/util/worker"
)

// Timeout wraps a worker with a given timeout
func Timeout(worker worker.Interface, t time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	run := func() chan error {
		ch := make(chan error)
		go func() {
			worker.Work(ctx, ch)
		}()
		return ch
	}
	select {
	case err := <-run():
		return err
	case <-ctx.Done():
		return errors.New("timeout exceeded")
	}
}
