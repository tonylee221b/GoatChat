package testutils

import "context"

type StubTx struct{}

func (s StubTx) WithinTx(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	return fn(ctx)
}
