package out

import (
	"context"
	"errors"
	"log/slog"

	identitysqlc "GoatChat/GoatChat/internal/identity/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/identity/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPgRepository struct {
	pool *pgxpool.Pool
}

func NewUserPgRepository(pool *pgxpool.Pool) *UserPgRepository {
	return &UserPgRepository{pool: pool}
}

func (u *UserPgRepository) Save(ctx context.Context, user domain.User) error {
	q := u.queries(ctx)

	_, err := q.CreateUser(ctx, identitysqlc.CreateUserParams{
		ID:          user.ID,
		Email:       &user.Email.Value,
		Username:    user.Username.Value,
		PhoneNumber: &user.PhoneNumber.Value,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPgRepository) FindByUsername(ctx context.Context, un domain.Username) (*domain.User, error) {
	q := u.queries(ctx)

	userFromDB, err := q.GetUserByUsername(ctx, identitysqlc.GetUserByUsernameParams{
		Username: un.Value,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New(domain.ErrUserNotFound)
		}

		slog.Error("something went wrong with db", "error", err)
		return nil, errors.New(domain.ErrDB)
	}
	slog.Debug("user found from db", "user", userFromDB.Username)

	username, err := domain.NewUsername(userFromDB.Username)
	if err != nil {
		return nil, err
	}

	return domain.NewUser(username)
}

func (u *UserPgRepository) queries(ctx context.Context) *identitysqlc.Queries {
	return identitysqlc.New(u.db(ctx))
}

func (u *UserPgRepository) db(ctx context.Context) identitysqlc.DBTX {
	if pgxTx := dbtx.FromContext(ctx); pgxTx != nil {
		return pgxTx
	}

	return u.pool
}
