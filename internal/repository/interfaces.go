package repository

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	user "github.com/alanzhumalin/bank/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User) error
	UserExists(ctx context.Context, phoneNumber string) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user domain.User) error
	GetByPhone(ctx context.Context, phone string) (domain.User, error)
}
