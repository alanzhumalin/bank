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
type currency struct {
	id         int
	name       string
	code       string
	symbol     string
	created_at time.Time
}

func NewCurrency(name string, code string, symbol string) currency {
	return currency{
		name:   name,
		code:   code,
		symbol: symbol,
	}
}

func NewCurrencyFromDB(id int, name string, code string, symbol string, created_at time.Time) currency {
	return currency{
		id:         id,
		name:       name,
		code:       code,
		symbol:     symbol,
		created_at: created_at,
	}
}
