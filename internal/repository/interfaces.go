package repository

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	user "github.com/alanzhumalin/bank/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User) error
	UserExists(ctx context.Context, phoneNumber string) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user domain.User) error
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
}

type CurrencyRepository interface {
	Create(ctx context.Context, c domain.Сurrency) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (domain.Сurrency, error)
	UpdateById(ctx context.Context, id int, name string, code string, symbol string) error
	Exists(ctx context.Context, code string) (bool, error)
	GetAll(ctx context.Context) ([]domain.Сurrency, error)
}

type TransferRepository interface {
	Create(ctx context.Context, t domain.Transfer) error
	GetById(ctx context.Context, id int) (domain.Transfer, error)
	GetAll(ctx context.Context) ([]domain.Transfer, error)
}

type AccountRepository interface {
	Create(ctx context.Context, a domain.Account) error
}
