package dto

import "errors"

var (
	ErrorFirstNameRequired = errors.New("First name is required")
	ErrorLastNameRequired  = errors.New("Last name is required")
	ErrorPhoneNumRequired  = errors.New("Phone number is required")
	ErrorPasswordRequired  = errors.New("Password is required")
	ErrorPasswordTooShort  = errors.New("Password must be at least 8 characters")
)

var (
	ErrorNameRequired   = errors.New("Name is required")
	ErrorCodeRequired   = errors.New("Code is required")
	ErrorSymbolRequired = errors.New("Symbol is required")
)
