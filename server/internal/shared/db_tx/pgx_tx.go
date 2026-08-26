package dbtx

import (
	"context"
	"log/slog"

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
	if tx := FromContext(ctx); tx != nil {
		slog.Debug("tx already exists")
		return fn(ctx)
	}

	pgxTx, err := p.pool.Begin(ctx)
	if err != nil {
		slog.Error("db cannot be connected", "error", err)
		return err
	}

	defer func() {
		slog.Debug("tx rolled back")
		_ = pgxTx.Rollback(ctx)
	}()

	txCtx := context.WithValue(ctx, pgxTxKey{}, pgxTx)

	if err := fn(txCtx); err != nil {
		slog.Debug("tx roll back triggered")
		return err // Rollback triggered
	}

	slog.Debug("tx committed")
	return pgxTx.Commit(ctx)
}

func FromContext(ctx context.Context) pgx.Tx {
	v, _ := ctx.Value(pgxTxKey{}).(pgx.Tx)
	return v
}
