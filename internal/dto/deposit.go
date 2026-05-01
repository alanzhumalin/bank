package dto

type CreateDepositRequest struct {
	AccountId int    `json:"account_id"`
	Amount    int    `json:"amount"`
	Source    string `json:"source"`
}
