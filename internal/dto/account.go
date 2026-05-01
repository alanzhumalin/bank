package dto

import (
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	UserId     int `json:"user_id"`
	CurrencyId int `json:"currency_id"`
}

type GetAccountResponse struct {
	Id         int             `json:"id"`
	UserId     int             `json:"user_id"`
	CurrencyId int             `json:"currency_id"`
	Balance    decimal.Decimal `json:"balance"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
}

func ToGetAccountResponse(a domain.Account) GetAccountResponse {
	return GetAccountResponse{
		Id:         a.Id,
		UserId:     a.UserId,
		CurrencyId: a.CurrencyId,
		Balance:    a.Balance,
		IsActive:   a.IsActive,
		CreatedAt:  a.CreatedAt,
	}
}

func (c *CreateAccountRequest) Validate() error {
	if c.UserId <= 0 {
		return ErrorUserIdRequired
	}

	if c.CurrencyId <= 0 {
		return ErrorCurrencyIdRequired
	}

	return nil
}

func (c *CreateAccountRequest) ToDomainModel() domain.Account {
	return domain.Account{
		UserId:     c.UserId,
		CurrencyId: c.CurrencyId,
	}
}
