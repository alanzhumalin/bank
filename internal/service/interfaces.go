package service

import (
	"context"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
)

type UserService interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (int, string, error)
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
	GetUserAccounts(ctx context.Context, userId int) ([]dto.GetAccountResponse, error)
}

type DepositService interface {
	Create(ctx context.Context, req dto.CreateDepositRequest, id int) error
}

type WithdrawalService interface {
	Create(ctx context.Context, req dto.CreateWindrawalRequest, userId int) (dto.IdempotencyResponse, error)
}

type TransactionService interface {
	GetByAccountId(ctx context.Context, accountId int, limit int, cursor string, currentUserId int) (dto.CursorResponse[dto.TransactionResponse], error)
	GetAll(ctx context.Context) ([]dto.TransactionResponse, error)
	GetByUserId(ctx context.Context, userId int, cursorValue string, limit int, currencies *[]string) (dto.CursorResponse[dto.TransactionResponse], error)
}

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest, ip string, device string) (*dto.TokenPair, error)
	Login(ctx context.Context, req dto.LoginRequest, ip string, device string) (*dto.TokenPair, error)
	UpdateSession(ctx context.Context, req dto.RefreshRequest) (*dto.TokenPair, string, error)
	LogoutFromAllDevices(ctx context.Context, userId int) error
	Logout(ctx context.Context, sessionId string, jti string, exp time.Time) error
}
