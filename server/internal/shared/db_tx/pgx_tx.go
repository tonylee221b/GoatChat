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
	if FromContext(ctx) != nil {
		slog.Debug("tx already exists", "tx_key", ctx.Value(pgxTxKey{}).(pgx.Tx))
		return fn(ctx)
	}

	pgxTx, err := p.pool.Begin(ctx)
	if err != nil {
		slog.Error("db cannot be connected", "error", err)
		return err
	}

	defer func() {
		slog.Debug("tx rolled back", "tx_key", ctx.Value(pgxTxKey{}).(pgx.Tx))
		_ = pgxTx.Rollback(ctx)
	}()

	txCtx := context.WithValue(ctx, pgxTxKey{}, pgxTx)

	if err := fn(txCtx); err != nil {
		slog.Debug("tx roll back triggered", "tx_key", ctx.Value(pgxTxKey{}).(pgx.Tx))
		return err // Rollback triggered
	}

	slog.Debug("tx committed", "tx_key", ctx.Value(pgxTxKey{}).(pgx.Tx))
	return pgxTx.Commit(ctx)
}

func FromContext(ctx context.Context) pgx.Tx {
	v, _ := ctx.Value(pgxTxKey{}).(pgx.Tx)
	return v
}
