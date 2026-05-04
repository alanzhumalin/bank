package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
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

func (ds *depositService) Create(ctx context.Context, req dto.CreateDepositRequest, id int) error {
	return ds.txManager.WithTx(ctx, func(ctx context.Context) error {

		acc, err := ds.accountRepo.GetByIdForUpdate(ctx, id)

		if err != nil {
			return err
		}

		if !acc.IsActive {
			return domain.AccountIsNotActive
		}

		mp, err := ds.transactionRepo.Create(ctx, domain.Transaction{
			Type:          "deposit",
			Amount:        req.Amount,
			AccountId:     id,
			Status:        "pending",
			StatusMessage: "transaction created",
		})

		if err != nil {
			return err
		}

		if err = ds.accountRepo.IncreaseBalance(ctx, req.Amount, id); err != nil {
			return err
		}

		if err = ds.depositRepo.Create(ctx, domain.Deposit{
			TransactionId: mp[id],
			AccountId:     id,
			Amount:        req.Amount,
			Source:        req.Source,
		}); err != nil {
			return err
		}

		if err = ds.transactionRepo.MarkTransaction(ctx, "completed", "transaction completed", mp[id]); err != nil {
			return err
		}

		return nil
	})
}
