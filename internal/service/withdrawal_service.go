package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type withdrawalService struct {
	withdrawalRepo  repository.WithdrawalRepository
	txManager       repository.TxManagerRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	idempotencyRepo repository.IdempotencyRepository
}

func NewWithdrawalService(idempotencyRepo repository.IdempotencyRepository, withdrawalRepo repository.WithdrawalRepository, txManager repository.TxManagerRepository, accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) WithdrawalService {
	return &withdrawalService{idempotencyRepo: idempotencyRepo, withdrawalRepo: withdrawalRepo, txManager: txManager, accountRepo: accountRepo, transactionRepo: transactionRepo}
}

func (w *withdrawalService) Create(ctx context.Context, req dto.CreateWindrawalRequest, userId int) (dto.IdempotencyResponse, error) {
	var idem dto.IdempotencyResponse

	err := w.txManager.WithTx(ctx, func(ctx context.Context) error {
		err := w.idempotencyRepo.Start(ctx, domain.Idempotency{
			UserId:         userId,
			IdempotencyKey: req.IdempotencyKey,
			Operation:      "withdraw",
		})

		if errors.Is(err, domain.ErrorIdempotencyAlreadyExists) {
			idempotency, err := w.idempotencyRepo.GetByKey(ctx, req.IdempotencyKey, userId)

			if err != nil {
				return err
			}

			switch idempotency.Status {
			case "completed":
				idem = dto.IdempotencyResponse{
					Status:   idempotency.Status,
					Response: json.RawMessage(idempotency.Response),
				}
				return nil
			case "pending":
				return domain.ErrorIdempotencyPending

			case "failed":
				return domain.ErrorIdempotencyFailed
			}

		}

		if err != nil {
			return err
		}

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
			CurrencyId:    account.CurrencyId,
			StatusMessage: "Withdraw transaction started",
		})
		if err != nil {
			return fmt.Errorf("error here: %w", err)
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

		status := "completed"
		response := "Withdraw transaction completed"
		transactionId := mp[req.AccountId]

		if err = w.transactionRepo.MarkTransaction(ctx, status, response, mp[req.AccountId]); err != nil {
			return err
		}

		mpa := map[string]any{
			"data": response,
		}

		responseByte, err := json.Marshal(mpa)

		if err != nil {
			return err
		}

		updatedAt := time.Now()

		// idempotency.TransactionId, idempotency.Status, idempotency.Response, idempotency.UpdatedAt, idempotency.IdempotencyKey,

		if err = w.idempotencyRepo.Complete(ctx, domain.Idempotency{
			UserId:         userId,
			TransactionId:  &transactionId,
			Status:         status,
			Response:       responseByte,
			UpdatedAt:      &updatedAt,
			IdempotencyKey: req.IdempotencyKey,
		}); err != nil {
			return err
		}

		idem = dto.IdempotencyResponse{
			Status:   status,
			Response: json.RawMessage(responseByte),
		}

		return nil

	})

	if err != nil {
		return dto.IdempotencyResponse{}, err
	}

	return idem, nil
}
