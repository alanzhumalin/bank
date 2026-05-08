package domain

import (
	"errors"
)

var (
	ErrorUserAlreadyExists = errors.New("User already exists")
	ErrorInvalidPassword   = errors.New("Invalid password")
	ErrorUserNotFound      = errors.New("User not found")
)

var (
	ErrorCurrencyAlreadyExists = errors.New("Currency already exists")
	ErrorCurrencyNotFound      = errors.New("Currency not found")
)

var (
	ErrorNotEnoughBalance     = errors.New("Not enough money")
	AccountNotFound           = errors.New("Account not found")
	AccountNotSupportCurrency = errors.New("Account does not support this currency")
	AccountIsNotActive        = errors.New("Account is not active")
)

var (
	ErrorTransferNotFound = errors.New("Transfer not found")
)

var (
	ErrorTransactionNotFound = errors.New("Transaction not found")
	ErrorPasswordNotCorrect  = errors.New("Password is not correct")
)
var (
	ErrorSessionNotFound = errors.New("Session not found")
)
