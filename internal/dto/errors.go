package dto

import "errors"

var (
	ErrorFirstNameRequired = errors.New("First name is required")
	ErrorLastNameRequired  = errors.New("Last name is required")
	ErrorPhoneNumRequired  = errors.New("Phone number is required")
	ErrorPasswordRequired  = errors.New("Password is required")
	ErrorPasswordTooShort  = errors.New("Password must be at least 8 characters")
	ErrorBirthdayRequired  = errors.New("Birthday is required")
)

var (
	ErrorNameRequired   = errors.New("Name is required")
	ErrorCodeRequired   = errors.New("Code is required")
	ErrorSymbolRequired = errors.New("Symbol is required")
)

var (
	ErrorSenderIdRequired   = errors.New("Sender account id is required")
	ErrorReceiverIdRequired = errors.New("Receiver account id is required")
	ErrorCurrencyIdRequired = errors.New("Currency id is required")
	ErrorAmountRequired     = errors.New("Amount is required")
)

var (
	ErrorUserIdRequired = errors.New("User id is required")
	CurrencyIdRequired  = errors.New("Currency id is required")
)

var (
	ErrorTransactionIdRequired = errors.New("Transaction id is required")
	ErrorAccountIdRequired     = errors.New("Account id is required")
	ErrorSourceRequired        = errors.New("Source is required")
)
var (
	SessionIdNotFound = errors.New("Session id not found")
	UserIdNotFound    = errors.New("User id not found")
	RoleNotFound      = errors.New("Role not found")
)

var (
	ErrorAmountIsRequired        = errors.New("Deposit amount is required")
	ErrorDepositSourceIsRequired = errors.New("Deposit source is required")
)
