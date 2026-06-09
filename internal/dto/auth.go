package dto

import "time"

type RegisterRequest struct {
	FirstName   string     `json:"firstname"`
	LastName    string     `json:"lastname"`
	Birthday    *time.Time `json:"birthday"`
	PhoneNumber string     `json:"phone_number"`
	Password    string     `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if r.PhoneNumber == "" {
		return ErrorPhoneNumRequired
	}

	if r.Password == "" {
		return ErrorPasswordRequired
	}

	if r.Birthday == nil {
		return ErrorBirthdayRequired
	}

	if r.FirstName == "" {
		return ErrorFirstNameRequired
	}
	if r.LastName == "" {
		return ErrorLastNameRequired
	}
	return nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (l *LoginRequest) Validate() error {
	if l.PhoneNumber == "" {
		return ErrorPhoneNumRequired
	}
	if l.Password == "" {
		return ErrorPasswordRequired
	}

	return nil
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserKey struct{}

type RoleKey struct{}

type SessionKey struct{}

type JTIKey struct {
}
type ExpKey struct{}

type OTPRequest struct {
	ChallengeId string `json:"challenge_id"`
	Code        string `json:"code"`
}

func (o *OTPRequest) Validate() error {
	if o.ChallengeId == "" {
		return ChallengeIdRequired
	}

	if o.Code == "" {
		return CodeRequired
	}

	return nil
}

type LoginResponse struct {
	ChallengeId string `json:"challenge_id"`
}
