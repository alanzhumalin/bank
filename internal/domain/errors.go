package domain

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("User already exists")
	ErrorInvalidPassword   = errors.New("Invalid password")
	ErrorUserNotFound      = errors.New("User not found")
)

var (
	ErrorCurrencyAlreadyExists = errors.New("Currency already exists")
	ErrorCurrencyNotFound      = errors.New("Currency not found")
)
