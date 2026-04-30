package dto

type CreateAccountRequest struct {
	UserId     int `json:"user_id"`
	CurrencyId int `json:"currency_id"`
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
