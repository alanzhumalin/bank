package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

type authRepository struct {
	pool *pgxpool.Pool
}

type LoginDetails struct {
	UserId   int
	Role     string
	Password string
}

func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &authRepository{pool: pool}
}

func (a *authRepository) GetDetails(context context.Context, phoneNumber string) (LoginDetails, error) {
	ctx, span := otel.Tracer("bank-api").Start(context, "AuthRepository.GetDetails")

	defer span.End()
	var loginDetails LoginDetails

	err := a.pool.QueryRow(ctx, `select id, role, password from users where phone_number = $1`, phoneNumber).Scan(&loginDetails.UserId, &loginDetails.Role, &loginDetails.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		return LoginDetails{}, domain.ErrorUserNotFound
	}

	return loginDetails, nil
}

func (a *authRepository) Сreate(ctx context.Context, session domain.Session) error {
	q := querier(ctx, a.pool)
	_, err := q.Exec(ctx, `insert into sessions(id, hashed_refresh_token, user_id, device, ip, created_at, expires_at) values($1, $2, $3, $4, $5, $6, $7)`, session.Id, session.HashedRefreshToken, session.UserId, session.Device, session.Ip, session.CreatedAt, session.ExpiresAt)

	if err != nil {
		return fmt.Errorf("error in creating session: %w", err)
	}

	return nil
}

func (a *authRepository) GetSessionById(ctx context.Context, sessionId string) (domain.Session, error) {

	q := querier(ctx, a.pool)

	var session domain.Session

	err := q.QueryRow(ctx, `select is_active, hashed_refresh_token, expires_at from sessions where id = $1 for update`, sessionId).Scan(&session.IsActive, &session.HashedRefreshToken, &session.ExpiresAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Session{}, domain.ErrorSessionNotFound
	}

	if err != nil {
		return domain.Session{}, fmt.Errorf("Error in getting session by id: %w", err)
	}

	return session, nil
}

func (a *authRepository) Update(ctx context.Context, newHashedToken string, expires_at time.Time, sessionId string) error {
	q := querier(ctx, a.pool)

	commangTag, err := q.Exec(ctx, `update sessions set hashed_refresh_token = $1, expires_at = $2 where id = $3`, newHashedToken, expires_at, sessionId)

	if err != nil {
		return fmt.Errorf("Error in updating a existing session: %w", err)
	}

	if commangTag.RowsAffected() == 0 {
		return domain.ErrorSessionNotFound
	}

	return nil
}

func (a *authRepository) Revoke(ctx context.Context, sessionId string) error {
	q := querier(ctx, a.pool)

	tag, err := q.Exec(ctx, `update sessions set is_active = false where id = $1`, sessionId)

	if err != nil {
		return fmt.Errorf("error in revoking the session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrorSessionNotFound
	}

	return nil

}

func (a *authRepository) RevokeAllUserDevices(ctx context.Context, userId int) error {
	q := querier(ctx, a.pool)

	tag, err := q.Exec(ctx, `update sessions set is_active = false where user_id = $1`, userId)

	if err != nil {
		return fmt.Errorf("error in revoking the session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrorSessionNotFound
	}

	return nil
}
