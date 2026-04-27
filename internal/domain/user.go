package domain

import "time"

// CREATE TABLE IF NOT EXISTS USERS (
//     id BIGINT generated always as identity primary key,
//     firstname text not null,
//     lastname text not null,
//     birthday TIMESTAMPtz not null,
//     phone_number text not null,
//     password text not null,
//     created_at TIMESTAMPtz not null default now()
//	   role text not null default 'user'
// );

type User struct {
	Id          int
	FirstName   string
	LastName    string
	Birthday    time.Time
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	Role        string
}
