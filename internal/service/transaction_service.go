package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (ts *transactionService) GetByAccountId(ctx context.Context, id int) ([]dto.TransactionResponse, error) {
	trs, err := ts.repo.GetByAccountId(ctx, id)

	if err != nil {
		return []dto.TransactionResponse{}, err
	}

	sl := make([]dto.TransactionResponse, 0, len(trs))

	for _, val := range trs {
		sl = append(sl, dto.ToTransactionResponse(val))
	}

	return sl, nil
}

func (ts *transactionService) GetAll(ctx context.Context) ([]dto.TransactionResponse, error) {
	trs, err := ts.repo.GetAll(ctx)

	if err != nil {
		return []dto.TransactionResponse{}, err
	}

	sl := make([]dto.TransactionResponse, 0, len(trs))

	for _, val := range trs {
		sl = append(sl, dto.ToTransactionResponse(val))
	}

	return sl, nil
}
