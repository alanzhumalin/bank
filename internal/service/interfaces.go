package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user domain.User) error
	GetByPhone(ctx context.Context, phone string) (dto.GetUserByPhoneResponse, error)
}
