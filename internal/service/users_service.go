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

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) error {

	err := s.repo.UserExists(ctx, req.PhoneNumber)
	if err != nil {

		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error occured by hashing the password %w", err)
	}
	u := domain.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Birthday:    req.Birthday,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
	}

	err = s.repo.Create(ctx, u)

	if err != nil {
		return fmt.Errorf("Error occured creating the user: %w", err)
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("Error delete user by id: %w", err)
	}
	return nil
}

func (s *userService) Update(ctx context.Context, user domain.User) error {
	err := s.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("Error update user: %w", err)
	}
	return nil
}

func (s *userService) GetByPhone(ctx context.Context, phone string) (dto.GetUserByPhoneResponse, error) {
	res, err := s.repo.GetByPhone(ctx, phone)

	if err != nil {
		return dto.GetUserByPhoneResponse{}, fmt.Errorf("Error get by phone: %w", err)
	}

	user := dto.GetUserByPhoneResponse{
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
