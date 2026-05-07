package service

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository, log zerolog.Logger) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetAll(ctx context.Context) ([]dto.GetUser, error) {
	users, err := s.repo.GetAll(ctx)

	if err != nil {
		return []dto.GetUser{}, fmt.Errorf("Error in get all, user_service: %w", err)
	}
	sl := make([]dto.GetUser, 0, len(users))

	for _, val := range users {
		sl = append(sl, dto.ToGetUser(val))
	}

	return sl, nil

}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (int, error) {

	err := s.repo.UserExists(ctx, req.PhoneNumber)
	if err != nil {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("Error occured by hashing the password %w", err)
	}
	u := domain.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Birthday:    req.Birthday,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
	}

	id, err := s.repo.Create(ctx, u)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Update(ctx context.Context, user domain.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) GetByPhone(ctx context.Context, phone string) (dto.GetUser, error) {
	res, err := s.repo.GetByPhone(ctx, phone)

	if err != nil {
		return dto.GetUser{}, err
	}

	user := dto.GetUser{
		Id:          res.Id,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		Birthday:    res.Birthday,
		PhoneNumber: res.PhoneNumber,
		CreatedAt:   res.CreatedAt,
		Role:        res.Role,
	}

	return user, nil
}
