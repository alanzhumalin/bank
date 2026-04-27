package dto

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
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
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
