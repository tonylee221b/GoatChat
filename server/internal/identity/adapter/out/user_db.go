package out

import (
	"context"
	"errors"
	"log/slog"
	"time"

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

func (u *UserPgRepository) ExistsByUsername(ctx context.Context, username domain.Username) (bool, error) {
	q := u.queries(ctx)

	isExist, err := q.ExistsUserByUsername(ctx, identitysqlc.ExistsUserByUsernameParams{
		Username: username.Value,
	})
	if err != nil {
		slog.Error("db error", "error", err)
		return false, errors.New(domain.ErrDB)
	}

	slog.Debug("user exists by username", "isUserExist", isExist)
	return isExist, nil
}

func (u *UserPgRepository) Save(ctx context.Context, user domain.User) error {
	q := u.queries(ctx)

	_, err := q.CreateUser(ctx, identitysqlc.CreateUserParams{
		ID:          user.ID,
		Username:    user.Username.Value,
		Email:       &user.Contact.Email.Value,
		PhoneNumber: &user.Contact.PhoneNumber.Value,
	})
	if err != nil {
		slog.Info("failed to create user", "user", user)
		return errors.New(domain.ErrDB)
	}

	slog.Debug("new user created", "user", user)
	return nil
}

func (u *UserPgRepository) Update(ctx context.Context, user domain.User) error {
	q := u.queries(ctx)

	_, err := q.UpdateContact(ctx, identitysqlc.UpdateContactParams{
		ID:          user.ID,
		PhoneNumber: &user.Contact.PhoneNumber.Value,
		Email:       &user.Contact.Email.Value,
	})
	if err != nil {
		slog.Info("failed to update user", "user", user)
		return errors.New(domain.ErrDB)
	}

	slog.Debug("user updated", "user", user)
	return nil
}

func (u *UserPgRepository) FindByUsername(ctx context.Context, un domain.Username) (*domain.User, error) {
	q := u.queries(ctx)

	userFromDB, err := q.GetUserByUsername(ctx, identitysqlc.GetUserByUsernameParams{
		Username: un.Value,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info("user not found", "username", un.Value)
			return nil, errors.New(domain.ErrUserNotFound)
		}

		slog.Error("something went wrong with db", "error", err)
		return nil, errors.New(domain.ErrDB)
	}

	slog.Debug("user found from db", "user", userFromDB.Username)
	return toDomain(userFromDB), nil
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

func toDomain(ufd identitysqlc.User) *domain.User {
	un, _ := domain.NewUsername(ufd.Username)
	c, _ := domain.NewContact(*ufd.PhoneNumber, *ufd.Email)

	var deletedAt *time.Time
	if !ufd.DeletedAt.Valid {
		deletedAt = nil
	}
	a := domain.NewAudit(ufd.CreatedAt.Time, ufd.UpdatedAt.Time, deletedAt)

	return &domain.User{
		ID:       ufd.ID,
		Username: un,
		Contact:  c,
		Audit:    a,
	}
}
