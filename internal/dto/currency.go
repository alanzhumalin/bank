package dto

import (
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
)

type CreateNewCurrencyRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

type UpdateCurrency struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

type GetCurrencyResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Symbol    string    `json:"symbol"`
	CreatedAt time.Time `json:"created_at"`
}

func NewGetCurrencyResponse(c domain.Сurrency) GetCurrencyResponse {
	return GetCurrencyResponse{
		Id:        c.Id,
		Name:      c.Name,
		Code:      c.Code,
		Symbol:    c.Symbol,
		CreatedAt: c.CreatedAt,
	}
}

func (c *CreateNewCurrencyRequest) Validate() error {
	if c.Name == "" {
		return ErrorNameRequired
	}

	if c.Code == "" {
		return ErrorCodeRequired
	}

	if c.Symbol == "" {
		return ErrorSymbolRequired
	}
	return nil
}

/*
CREATE TABLE IF NOT EXISTS CURRENCIES(

	id bigint generated always as identity primary key,
	name text not null,
	code char(3) not null unique,
	symbol VARCHAR(5) not null,
	created_at TIMESTAMPtz not null DEFAULT now()

);
*/
