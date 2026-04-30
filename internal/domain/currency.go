package domain

import "time"

type Сurrency struct {
	Id        int
	Name      string
	Code      string
	Symbol    string
	CreatedAt time.Time
}
