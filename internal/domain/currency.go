package domain

import "time"

/*
CREATE TABLE IF NOT EXISTS CURRENCIES(

	id bigint generated always as identity primary key,
	name text not null,
	code char(3) not null unique,
	symbol VARCHAR(5) not null,
	created_at TIMESTAMPtz not null DEFAULT now()

);
*/
type Сurrency struct {
	Id         int
	Name       string
	Code       string
	Symbol     string
	Created_at time.Time
}

func NewCurrency(name string, code string, symbol string) Сurrency {
	return Сurrency{
		Name:   name,
		Code:   code,
		Symbol: symbol,
	}
}

func NewCurrencyFromDB(id int, name string, code string, symbol string, created_at time.Time) Сurrency {
	return Сurrency{
		Id:         id,
		Name:       name,
		Code:       code,
		Symbol:     symbol,
		Created_at: created_at,
	}
}
