package dbtx

import "context"

type Tx interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
