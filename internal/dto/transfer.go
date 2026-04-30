package dto

type CreateTransferRequest struct {
	SenderAccountId   int `json:"sender_account_id"`
	ReceiverAccountId int `json:"receiver_account_id"`
	CurrencyId        int `json:"currency_id"`
	Amount            int `json:"amount"`
}

func (c *CreateTransferRequest) Validate() error {
	if c.SenderAccountId <= 0 {
		return ErrorSenderIdRequired
	}

	if c.ReceiverAccountId <= 0 {
		return ErrorReceiverIdRequired
	}

	if c.CurrencyId <= 0 {
		return ErrorCurrencyIdRequired
	}

	if c.Amount <= 0 {
		return ErrorAmountRequired
	}

	return nil
}
