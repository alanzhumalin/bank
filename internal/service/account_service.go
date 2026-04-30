package service

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (a *accountService) Create(ctx context.Context, req dto.CreateAccountRequest) error {

	err := a.repo.Create(ctx, req.ToDomainModel())
	if err != nil {
		return fmt.Errorf("Error in creating account: %w", err)
	}

	return nil
}

func (a *accountService) DeleteById(ctx context.Context, id int) error {
	err := a.repo.DeleteById(ctx, id)

	if err != nil {
		return fmt.Errorf("Error in deleting account by id: %w", err)
	}

	return err
}

func (a *accountService) GetAll(ctx context.Context) ([]dto.GetAccountResponse, error) {
	accounts, err := a.repo.GetAll(ctx)
	if err != nil {
		return []dto.GetAccountResponse{}, fmt.Errorf("Error in get all account: %w", err)
	}

	sl := make([]dto.GetAccountResponse, 0, len(accounts))

	for _, val := range accounts {

		sl = append(sl, dto.ToGetAccountResponse(val))
	}

	return sl, nil
}
