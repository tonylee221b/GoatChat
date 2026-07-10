package dbtx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxTxKey struct{}

type pgxTx struct {
	pool *pgxpool.Pool
}

func NewPgxTx(pool *pgxpool.Pool) *pgxTx {
	return &pgxTx{pool}
}

func (p *pgxTx) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if FromContext(ctx) != nil {
		return fn(ctx)
	}

	pgxTx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = pgxTx.Rollback(ctx)
	}()

	txCtx := context.WithValue(ctx, pgxTxKey{}, pgxTx)

	if err := fn(txCtx); err != nil {
		return err // Rollback triggered
	}

	return pgxTx.Commit(ctx)
}

func FromContext(ctx context.Context) pgx.Tx {
	v, _ := ctx.Value(pgxTxKey{}).(pgx.Tx)
	return v
}
