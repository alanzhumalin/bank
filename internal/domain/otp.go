package domain

type OTPDetail struct {
	UserId      int    `json:"user_id"`
	Attempt     int    `json:"attempt"`
	PhoneNumber string `json:"phone_number"`
	CodeHash    string `json:"code_hash"`
}
