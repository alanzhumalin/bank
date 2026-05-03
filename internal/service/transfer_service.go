package service

import (
	"context"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type transferService struct {
	transferRepo    repository.TransferRepository
	txManager       repository.TxManagerRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewTransferService(transferRepo repository.TransferRepository, txManager repository.TxManagerRepository, accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) TransferService {
	return &transferService{transferRepo: transferRepo, txManager: txManager, accountRepo: accountRepo, transactionRepo: transactionRepo}
}

func (t *transferService) Create(ctx context.Context, req dto.CreateTransferRequest) error {

	err := t.txManager.WithTx(ctx, func(ctx context.Context) error {
		acc1, acc2, err := t.accountRepo.SelectTwoAccountsForUpdate(ctx, req.SenderAccountId, req.ReceiverAccountId)

		if err != nil {
			return err
		}

		if acc2.CurrencyId != req.CurrencyId {
			return domain.AccountNotSupportCurrency
		}

		if !acc2.IsActive {
			return domain.AccountIsNotActive
		}

		if acc1.CurrencyId != req.CurrencyId {
			return domain.AccountNotSupportCurrency
		}
		if !acc1.IsActive {
			return domain.AccountIsNotActive
		}

		if acc1.Balance.LessThan(req.Amount) {
			return domain.ErrorNotEnoughBalance
		}

		trans := domain.Transaction{
			Type:          "transfer",
			Amount:        req.Amount,
			AccountId:     req.SenderAccountId,
			Status:        "pending",
			StatusMessage: "transaction started",
		}

		id, err := t.transactionRepo.Create(ctx, trans)

		if err != nil {
			return err
		}

		err = t.accountRepo.DecreaseBalance(ctx, req.Amount, acc1.Id)

		if err != nil {

			return err
		}

		err = t.accountRepo.IncreaseBalance(ctx, req.Amount, acc2.Id)

		if err != nil {
			return err
		}

		err = t.transactionRepo.MarkTransaction(ctx, "completed", "transaction successfuly completed", id)

		if err != nil {
			return err
		}

		return nil

	})

	return err
}
