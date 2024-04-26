package api

import (
	"context"
	"time"
)

type noTimeoutContext struct {
	inner context.Context
}

func (no noTimeoutContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (no noTimeoutContext) Done() <-chan struct{} {
	return nil
}

func (no noTimeoutContext) Err() error {
	return nil
}

func (no noTimeoutContext) Value(key any) any {
	return no.inner.Value(key)
}

func noCancelContext(ctx context.Context) context.Context {
	return noTimeoutContext{
		inner: ctx,
	}
}
