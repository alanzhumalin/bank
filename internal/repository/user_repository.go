package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	user "github.com/alanzhumalin/bank/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{
		pool: pool,
	}
}

// CREATE TABLE IF NOT EXISTS USERS (
//     id BIGINT generated always as identity primary key,
//     firstname text not null,
//     lastname text not null,
//     birthday TIMESTAMPtz not null,
//     phone_number text not null,
//     password text not null,
//     created_at TIMESTAMPtz not null default now()
//	   role text not null default 'user'
// );

type User struct {
	Id          int
	FirstName   string
	LastName    string
	Birthday    time.Time
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	Role        string
}

func (u *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	rows, err := u.pool.Query(ctx, `select id, firstname, lastname, birthday, phone_number, created_at, role from users`)

	if err != nil {
		return []domain.User{}, fmt.Errorf("Error in get all, user_repository: %w", err)
	}

	sl := make([]domain.User, 0)

	for rows.Next() {
		var req domain.User

		err := rows.Scan(&req.Id, &req.FirstName, &req.LastName, &req.Birthday, &req.PhoneNumber, &req.CreatedAt, &req.Role)

		if err != nil {
			return []domain.User{}, fmt.Errorf("Error in a loop, after scan: %w", err)
		}

		sl = append(sl, req)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return []domain.User{}, fmt.Errorf("Error in get all, in rows: %w", err)
	}
	return sl, nil
}

func (u *userRepository) UserExists(ctx context.Context, phoneNumber string) error {
	var num int
	err := u.pool.QueryRow(ctx, `select id from users where phone_number = $1`, phoneNumber).Scan(&num)

	if err == nil {
		return domain.ErrorUserAlreadyExists
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	return fmt.Errorf("user existence check: %w", err)
}

func (u *userRepository) Create(ctx context.Context, user user.User) (int, error) {
	var id int
	err := u.pool.QueryRow(ctx, `insert into users(firstname, lastname, birthday, phone_number, password) values($1,$2,$3,$4,$5) returning id`, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Password).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("user create: %w", err)
	}

	return id, nil
}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	res, err := u.pool.Exec(ctx, `delete from users where id = $1`, id)

	if err != nil {
		return fmt.Errorf("Error delete user by id: %w", err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrorUserNotFound
	}

	return nil
}

func (u *userRepository) Update(ctx context.Context, user domain.User) error {
	res, err := u.pool.Exec(ctx, `update users set firstname = $1, lastname = $2 where id = $3`, user.FirstName, user.LastName, user.Id)

	if err != nil {
		return fmt.Errorf("Error update user: %w", err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrorUserNotFound
	}

	return nil
}

func (u *userRepository) GetByPhone(ctx context.Context, phone string) (domain.User, error) {
	var user domain.User
	err := u.pool.QueryRow(ctx, `select * from users where phone_number = $1`, phone).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Birthday, &user.PhoneNumber, &user.Password, &user.CreatedAt, &user.Role)

	if err == nil {
		return user, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrorUserNotFound
	}

	return domain.User{}, fmt.Errorf("Error Get by phone: %w", err)
}
