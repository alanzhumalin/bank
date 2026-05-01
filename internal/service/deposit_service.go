package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type depositService struct {
	depositRepo     repository.DepositRepository
	accountRepo     repository.AccountRepository
	txManager       repository.TxManagerRepository
	transactionRepo repository.TransactionRepository
}

func NewDepositService(depoRepo repository.DepositRepository, accountRepo repository.AccountRepository, txManager repository.TxManagerRepository, transactionRepo repository.TransactionRepository) DepositService {
	return &depositService{
		depositRepo:     depoRepo,
		accountRepo:     accountRepo,
		txManager:       txManager,
		transactionRepo: transactionRepo,
	}
}

func (ds *depositService) Create(ctx context.Context, req dto.CreateDepositRequest) error {
	return ds.txManager.WithTx(ctx, func(ctx context.Context) error {

		return nil
	})
}
