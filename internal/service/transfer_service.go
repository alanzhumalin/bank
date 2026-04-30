package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type transferService struct {
	repo repository.TransferRepository
}

func NewTransferService(repo repository.TransferRepository) TransferService {
	return &transferService{repo: repo}
}

func (t *transferService) Create(ctx context.Context, req dto.CreateTransferRequest) error {
	transfer := domain.NewTransfer(req.SenderAccountId, req.ReceiverAccountId, req.CurrencyId, float64(req.Amount))

	return t.repo.Create(ctx, transfer)
}

func (t *transferService) GetAll(ctx context.Context) ([]dto.TransferResponse, error) {
	transfers, err := t.repo.GetAll(ctx)

	if err != nil {
		return []dto.TransferResponse{}, err
	}

	sl := make([]dto.TransferResponse, 0, len(transfers))

	for _, val := range transfers {
		sl = append(sl, *dto.ToTransferResponse(val))
	}

	return sl, nil
}

func (t *transferService) GetById(ctx context.Context, id int) (dto.TransferResponse, error) {
	transfer, err := t.repo.GetById(ctx, id)

	if err != nil {
		return dto.TransferResponse{}, err
	}

	return *dto.ToTransferResponse(transfer), nil
}
