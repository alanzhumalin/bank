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

	return t.txManager.WithTx(ctx, func(ctx context.Context) error {
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

		transaction1 := domain.Transaction{
			Type:          "transfer",
			Amount:        req.Amount,
			AccountId:     req.SenderAccountId,
			CurrencyId:    req.CurrencyId,
			Status:        "pending",
			StatusMessage: "transaction started",
		}
		transaction2 := domain.Transaction{
			Type:          "transfer",
			Amount:        req.Amount,
			AccountId:     req.ReceiverAccountId,
			CurrencyId:    req.CurrencyId,
			Status:        "pending",
			StatusMessage: "transaction started",
		}

		mp, err := t.transactionRepo.Create(ctx, transaction1, transaction2)

		if err != nil {
			return err
		}

		if err = t.accountRepo.DecreaseBalance(ctx, req.Amount, acc1.Id); err != nil {
			return err
		}

		if err = t.accountRepo.IncreaseBalance(ctx, req.Amount, acc2.Id); err != nil {
			return err

		}

		transfer1 := domain.Transfer{
			TransactionId:     mp[req.SenderAccountId],
			SenderAccountId:   req.SenderAccountId,
			ReceiverAccountId: req.ReceiverAccountId,
			CurrencyId:        req.CurrencyId,
			Amount:            req.Amount,
		}

		transfer2 := domain.Transfer{
			TransactionId:     mp[req.ReceiverAccountId],
			SenderAccountId:   req.SenderAccountId,
			ReceiverAccountId: req.ReceiverAccountId,
			CurrencyId:        req.CurrencyId,
			Amount:            req.Amount,
		}

		if err = t.transferRepo.Create(ctx, transfer1, transfer2); err != nil {
			return err
		}

		for _, val := range mp {
			if err = t.transactionRepo.MarkTransaction(ctx, "completed", "transaction successfuly completed", val); err != nil {
				return err
			}
		}

		return nil

	})

}
