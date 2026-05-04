package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type withdrawalService struct {
	withdrawalRepo  repository.WithdrawalRepository
	txManager       repository.TxManagerRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewWithdrawalService(withdrawalRepo repository.WithdrawalRepository, txManager repository.TxManagerRepository, accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) WithdrawalService {
	return &withdrawalService{withdrawalRepo: withdrawalRepo, txManager: txManager, accountRepo: accountRepo, transactionRepo: transactionRepo}
}

func (w *withdrawalService) Create(ctx context.Context, req dto.CreateWindrawalRequest) error {
	return w.txManager.WithTx(ctx, func(ctx context.Context) error {
		account, err := w.accountRepo.GetByIdForUpdate(ctx, req.AccountId)

		if err != nil {
			return err
		}

		if !account.IsActive {
			return domain.AccountIsNotActive
		}

		if account.Balance.LessThan(req.Amount) {
			return domain.ErrorNotEnoughBalance
		}

		mp, err := w.transactionRepo.Create(ctx, domain.Transaction{
			Type:          "withdraw",
			Amount:        req.Amount,
			AccountId:     req.AccountId,
			StatusMessage: "Withdraw transaction started",
		})
		if err != nil {
			return err
		}

		if err = w.accountRepo.DecreaseBalance(ctx, req.Amount, req.AccountId); err != nil {
			return err
		}

		if err = w.withdrawalRepo.Create(ctx, domain.Withdrawal{
			TransactionId: mp[req.AccountId],
			AccountId:     req.AccountId,
			Amount:        req.Amount,
			Source:        req.Source,
		}); err != nil {
			return err
		}

		if err = w.transactionRepo.MarkTransaction(ctx, "completed", "Withdraw transaction completed", mp[req.AccountId]); err != nil {
			return err
		}

		return nil

	})
}
