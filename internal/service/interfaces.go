package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
)

type UserService interface {
	Create(ctx context.Context, req dto.CreateUserRequest) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user domain.User) error
	GetByPhone(ctx context.Context, phone string) (dto.GetUser, error)
	GetAll(ctx context.Context) ([]dto.GetUser, error)
}

type CurrencyService interface {
	Create(ctx context.Context, req dto.CreateNewCurrencyRequest) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, req dto.UpdateCurrency) error
	GetAll(ctx context.Context) ([]dto.GetCurrencyResponse, error)
	GetById(ctx context.Context, id int) (dto.GetCurrencyResponse, error)
}

type TransferService interface {
	Create(ctx context.Context, req dto.CreateTransferRequest) error
}

type AccountService interface {
	Create(ctx context.Context, req dto.CreateAccountRequest) error
	DeleteById(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]dto.GetAccountResponse, error)
}

type DepositService interface {
	Create(ctx context.Context, req dto.CreateDepositRequest, id int) error
}

type WithdrawalService interface {
	Create(ctx context.Context, req dto.CreateWindrawalRequest) error
}
