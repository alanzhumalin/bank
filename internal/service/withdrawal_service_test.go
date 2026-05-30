package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/shopspring/decimal"
)

type fakeWithdrawalRepo struct {
	createCalled bool
	createError  error
}

func (wr *fakeWithdrawalRepo) Create(ctx context.Context, w domain.Withdrawal) error {
	wr.createCalled = true

	if wr.createError != nil {
		return wr.createError
	}

	return nil
}

func TestWithdrawalCreateSuccess(t *testing.T) {
	withdrawRepo := &fakeWithdrawalRepo{}
	txManager := &fakeTxManager{}
	accountRepo := &fakeAccountRepo{
		account: domain.Account{
			Id:         1,
			UserId:     1,
			CurrencyId: 1,
			Balance:    decimal.NewFromInt(1000000),
			IsActive:   true,
			CreatedAt:  time.Now(),
		},
	}
	transactionRepo := &fakeTransactionRepo{}
	req := dto.CreateWindrawalRequest{
		AccountId: 1,
		Amount:    decimal.NewFromInt(100000),
		Source:    "terminal",
	}
	srv := NewWithdrawalService(withdrawRepo, txManager, accountRepo, transactionRepo)

	err := srv.Create(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !accountRepo.calledGetByIdForUpdate {
		t.Fatal("Expected calledGetByIdForUpdate to be called")
	}

	if !transactionRepo.calledCreate {
		t.Fatal("Expected transaction create to be called")
	}

	if !accountRepo.calledDecreaseBalance {
		t.Fatalf("Expected account decrease balance to be called")
	}

	if !withdrawRepo.createCalled {
		t.Fatalf("Expected withdrawal create to be called")
	}

	if !transactionRepo.calledMarkTransaction {
		t.Fatalf("Expected transaction marked to be called")
	}

}

func TestWithdrawalErrors(t *testing.T) {
	tests := []struct {
		name                        string
		wantedError                 error
		req                         dto.CreateWindrawalRequest
		wantTxManagerCalled         bool
		wantGetByIdForUpdateCalled  bool
		wantTransactionCreateCalled bool
		wantDecreaseBalanceCalled   bool
		wantWithdrawalCreateCalled  bool
		wantMarkTransactionCalled   bool
		setup                       func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo)
	}{
		{
			name:        "Error account is not active",
			wantedError: domain.AccountIsNotActive,
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(10000),
				Source:    "terminal",
			},
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(100),
					IsActive:   false,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: false,
			wantDecreaseBalanceCalled:   false,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name:        "Account hass less money",
			wantedError: domain.ErrorNotEnoughBalance,
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(10000),
				Source:    "terminal",
			},
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(100),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: false,
			wantDecreaseBalanceCalled:   false,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			withdrawalRepo := &fakeWithdrawalRepo{}
			txManager := &fakeTxManager{}
			accountRepo := &fakeAccountRepo{}
			transactionRepo := &fakeTransactionRepo{}

			test.setup(withdrawalRepo, txManager, accountRepo, transactionRepo)

			srv := NewWithdrawalService(withdrawalRepo, txManager, accountRepo, transactionRepo)

			err := srv.Create(context.Background(), test.req)

			if !errors.Is(err, test.wantedError) {
				t.Fatalf("Expected error: %v, got %v", test.wantedError, err)
			}

			if test.wantTxManagerCalled != txManager.calledTx {
				t.Fatalf("Expected wantTxManagerCalled: %v, got %v", test.wantTxManagerCalled, txManager.calledTx)
			}

			if test.wantGetByIdForUpdateCalled != accountRepo.calledGetByIdForUpdate {
				t.Fatalf("Expected calledGetByIdForUpdate: %v, got %v", test.wantGetByIdForUpdateCalled, accountRepo.calledGetByIdForUpdate)
			}

			if test.wantTransactionCreateCalled != transactionRepo.calledCreate {
				t.Fatalf("Expected wantTransactionCreateCalled: %v, got %v", test.wantTransactionCreateCalled, transactionRepo.calledCreate)
			}

			if test.wantDecreaseBalanceCalled != accountRepo.calledDecreaseBalance {
				t.Fatalf("Expected wantDecreaseBalanceCalled: %v, got %v", test.wantDecreaseBalanceCalled, accountRepo.calledDecreaseBalance)
			}

			if test.wantWithdrawalCreateCalled != withdrawalRepo.createCalled {
				t.Fatalf("Expected wantWithdrawalCreateCalled: %v, got %v", test.wantWithdrawalCreateCalled, withdrawalRepo.createCalled)
			}

			if test.wantMarkTransactionCalled != transactionRepo.calledMarkTransaction {
				t.Fatalf("Expected wantMarkTransactionCalled: %v, got %v", test.wantMarkTransactionCalled, transactionRepo.calledCreate)
			}

		})
	}
}

func TestWithdrawalRepoErrors(t *testing.T) {
	tests := []struct {
		name                        string
		req                         dto.CreateWindrawalRequest
		wantError                   error
		wantTxManagerCalled         bool
		wantGetByIdForUpdateCalled  bool
		wantTransactionCreateCalled bool
		wantDecreaseBalanceCalled   bool
		wantWithdrawalCreateCalled  bool
		wantMarkTransactionCalled   bool
		setup                       func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo)
	}{
		{
			name: "Error withtx db",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: ErrorDb,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				txManager.withTxError = ErrorDb
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  false,
			wantTransactionCreateCalled: false,
			wantDecreaseBalanceCalled:   false,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error getbyidforupdate db",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: ErrorDb,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.getByIdForUpdateError = ErrorDb
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: false,
			wantDecreaseBalanceCalled:   false,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error getbyidforupdate account not found",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: domain.AccountNotFound,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.getByIdForUpdateError = domain.AccountNotFound
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: false,
			wantDecreaseBalanceCalled:   false,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error transaction create db",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: ErrorDb,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				transactionRepo.createErr = ErrorDb
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantDecreaseBalanceCalled:   false,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error decrease balance db",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: ErrorDb,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.decreaseBalanceError = ErrorDb
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantDecreaseBalanceCalled:   true,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error decrease balance account not found",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: domain.AccountNotFound,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.decreaseBalanceError = domain.AccountNotFound
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantDecreaseBalanceCalled:   true,
			wantWithdrawalCreateCalled:  false,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error withdrawal create db",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: ErrorDb,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				withdrawalRepo.createError = ErrorDb
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantDecreaseBalanceCalled:   true,
			wantWithdrawalCreateCalled:  true,
			wantMarkTransactionCalled:   false,
		},

		{
			name: "Error mark transaction db",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: ErrorDb,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				transactionRepo.markTransactionErr = ErrorDb
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantDecreaseBalanceCalled:   true,
			wantWithdrawalCreateCalled:  true,
			wantMarkTransactionCalled:   true,
		},

		{
			name: "Error mark transaction not found",
			req: dto.CreateWindrawalRequest{
				AccountId: 1,
				Amount:    decimal.NewFromInt(1000),
				Source:    "terminal",
			},
			wantError: domain.ErrorTransactionNotFound,
			setup: func(withdrawalRepo *fakeWithdrawalRepo, txManager *fakeTxManager, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				transactionRepo.markTransactionErr = domain.ErrorTransactionNotFound
				accountRepo.account = domain.Account{
					Id:         1,
					UserId:     1,
					CurrencyId: 1,
					Balance:    decimal.NewFromInt(1000000),
					IsActive:   true,
					CreatedAt:  time.Now(),
				}
			},
			wantTxManagerCalled:         true,
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantDecreaseBalanceCalled:   true,
			wantWithdrawalCreateCalled:  true,
			wantMarkTransactionCalled:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			withdrawalRepo := &fakeWithdrawalRepo{}
			txManager := &fakeTxManager{}
			accountRepo := &fakeAccountRepo{}
			transactionRepo := &fakeTransactionRepo{}

			test.setup(withdrawalRepo, txManager, accountRepo, transactionRepo)

			srv := NewWithdrawalService(withdrawalRepo, txManager, accountRepo, transactionRepo)

			err := srv.Create(context.Background(), test.req)

			if !errors.Is(err, test.wantError) {
				t.Fatalf("Expected error: %v, got %v", test.wantError, err)
			}

			if test.wantTxManagerCalled != txManager.calledTx {
				t.Fatalf("Expected wantTxManagerCalled: %v, got %v", test.wantTxManagerCalled, txManager.calledTx)
			}

			if test.wantGetByIdForUpdateCalled != accountRepo.calledGetByIdForUpdate {
				t.Fatalf("Expected wantGetByIdForUpdateCalled: %v, got %v", test.wantGetByIdForUpdateCalled, accountRepo.calledGetByIdForUpdate)
			}

			if test.wantTransactionCreateCalled != transactionRepo.calledCreate {
				t.Fatalf("Expected wantTransactionCreateCalled: %v, got %v", test.wantTransactionCreateCalled, transactionRepo.calledCreate)
			}

			if test.wantDecreaseBalanceCalled != accountRepo.calledDecreaseBalance {
				t.Fatalf("Expected wantDecreaseBalanceCalled: %v, got %v", test.wantDecreaseBalanceCalled, accountRepo.calledDecreaseBalance)
			}

			if test.wantWithdrawalCreateCalled != withdrawalRepo.createCalled {
				t.Fatalf("Expected wantWithdrawalCreateCalled: %v, got %v", test.wantWithdrawalCreateCalled, withdrawalRepo.createCalled)
			}

			if test.wantMarkTransactionCalled != transactionRepo.calledMarkTransaction {
				t.Fatalf("Expected wantMarkTransactionCalled: %v, got %v", test.wantMarkTransactionCalled, transactionRepo.calledMarkTransaction)
			}
		})
	}
}
