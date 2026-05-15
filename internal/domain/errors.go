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
	AccountAlreadyExists      = errors.New("Account already exists")
)

var (
	ErrorTransferNotFound = errors.New("Transfer not found")
)

var (
	ErrorTransactionNotFound = errors.New("Transaction not found")
	ErrorPasswordNotCorrect  = errors.New("Password is not correct")
)
var (
	ErrorSessionNotFound       = errors.New("Session not found")
	ErrorIncorrectRefreshToken = errors.New("Incorrect refresh token")
	ErrorRefreshTokenExpired   = errors.New("Refresh token expired")
	ErrorSessionNotActive      = errors.New("Session is not active")
)

var (
	ErrorNoAccounts = errors.New("User has no accounts")
)
