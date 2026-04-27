package dto

import "errors"

var (
	ErrorFirstNameRequired = errors.New("First name is required")
	ErrorLastNameRequired  = errors.New("Last name is required")
	ErrorPhoneNumRequired  = errors.New("Phone number is required")
	ErrorPasswordRequired  = errors.New("Password is required")
	ErrorPasswordTooShort  = errors.New("Password must be at least 8 characters")
)
