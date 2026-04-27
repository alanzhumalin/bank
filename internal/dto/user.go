package dto

import (
	"time"
)

// CREATE TABLE IF NOT EXISTS USERS (
//
//	    id BIGINT generated always as identity primary key,
//	    firstname text not null,
//	    lastname text not null,
//	    birthday TIMESTAMPtz not null,
//	    phone_number text not null,
//	    password text not null,
//	    created_at TIMESTAMPtz not null default now()
//		   role text not null default 'user'
//
// );

type CreateUserRequest struct {
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Birthday    time.Time `json:"birthday"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
}

type GetUserByPhoneRequest struct {
	PhoneNumber string `json:"phone_number"`
}

func (g *GetUserByPhoneRequest) Validate() error {
	if g.PhoneNumber == "" {
		return ErrorPhoneNumRequired
	}
	return nil
}

type GetUserByPhoneResponse struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Birthday    time.Time `json:"birthday"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	Role        string    `json:"role"`
}

func (c *CreateUserRequest) Validate() error {
	if c.FirstName == "" {
		return ErrorFirstNameRequired
	}
	if c.LastName == "" {
		return ErrorLastNameRequired
	}
	if c.Password == "" {
		return ErrorPasswordRequired
	}

	if len(c.Password) < 8 {
		return ErrorPasswordTooShort
	}
	return nil
}
