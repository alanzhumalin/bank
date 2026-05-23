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

//  depositRepo     repository.DepositRepository
// 	accountRepo     repository.AccountRepository
// 	txManager       repository.TxManagerRepository
// 	transactionRepo repository.TransactionRepository

type fakeDepositRepo struct {
	createCalled bool
	createError  error
}

func (f *fakeDepositRepo) Create(ctx context.Context, d domain.Deposit) error {
	f.createCalled = true

	if f.createError != nil {
		return f.createError
	}

	return nil
}

func TestDepositService_Create_Success(t *testing.T) {

	depositRepo := &fakeDepositRepo{}
	accountRepo := &fakeAccountRepo{
		account: domain.Account{
			IsActive:   true,
			CurrencyId: 1,
		},
	}
	txManager := &fakeTxManager{}
	transactionRepo := &fakeTransactionRepo{}

	srv := NewDepositService(depositRepo, accountRepo, txManager, transactionRepo)

	err := srv.Create(context.Background(), dto.CreateDepositRequest{
		Amount: decimal.NewFromInt(1000),
		Source: "terminal",
	}, 1)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !txManager.calledTx {
		t.Fatal("Expected with tx to be called")
	}

	if !accountRepo.calledGetByIdForUpdate {
		t.Fatal("Expected calledGetByIdForUpdate to be called")
	}

	if !transactionRepo.calledCreate {
		t.Fatal("Expected transactionCreate to be called")
	}

	if !accountRepo.calledIncreaseBalance {
		t.Fatal("Expected increaseBalance to be called")
	}

	if !depositRepo.createCalled {
		t.Fatal("Expected createDeposit to be called")
	}

	if !transactionRepo.calledMarkTransaction {
		t.Fatal("Expected calledMarkTransaction to be called")
	}

}

func TestDepositCreateErrors(t *testing.T) {
	tests := []struct {
		name string

		id int

		req dto.CreateDepositRequest

		wantError error
	}{
		{
			name: "Account is not active",
			id:   1,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(1000),
				Source: "terminal",
			},
			wantError: domain.AccountIsNotActive,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			depositRepo := &fakeDepositRepo{}
			accountRepo := &fakeAccountRepo{}
			txManager := &fakeTxManager{}
			transactionRepo := &fakeTransactionRepo{}

			srv := NewDepositService(depositRepo, accountRepo, txManager, transactionRepo)

			err := srv.Create(context.Background(), test.req, test.id)

			if !errors.Is(err, test.wantError) {
				t.Fatalf("Expected error %v, got %v", test.wantError, err)
			}

			if !txManager.calledTx {
				t.Fatalf("Expected txmanager to be called")
			}

			if transactionRepo.calledCreate {
				t.Fatalf("Expected transaction create not be called")
			}

			if accountRepo.calledIncreaseBalance {
				t.Fatalf("Expected account increase balance not be called")
			}

			if depositRepo.createCalled {
				t.Fatalf("Expected deposit create not be called")
			}

			if transactionRepo.calledMarkTransaction {
				t.Fatalf("Expected mark transaction not be called")
			}

		})
	}

}

func TestDepositCreateRepoErrors(t *testing.T) {
	tests := []struct {
		name                        string
		id                          int
		wantedError                 error
		wantGetByIdForUpdateCalled  bool
		wantTransactionCreateCalled bool
		wantIncreaseBalance         bool
		wantDepositCreate           bool
		wantMarkTransactionCalled   bool
		req                         dto.CreateDepositRequest
		setup                       func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo)
	}{
		{
			id:          1,
			name:        "Error db from with tx, begin transaction",
			wantedError: ErrorDb,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				txRepo.withTxError = ErrorDb
			},
			wantGetByIdForUpdateCalled:  false,
			wantTransactionCreateCalled: false,
			wantIncreaseBalance:         false,
			wantDepositCreate:           false,
			wantMarkTransactionCalled:   false,
		},

		{
			id:          1,
			name:        "Error db from getbyidforupdate account",
			wantedError: ErrorDb,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: false,
			wantIncreaseBalance:         false,
			wantDepositCreate:           false,
			wantMarkTransactionCalled:   false,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.getByIdForUpdateError = ErrorDb
			},
		},

		{
			id:          1,
			name:        "Account not found from getbyidforupdate account",
			wantedError: domain.AccountNotFound,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: false,
			wantIncreaseBalance:         false,
			wantDepositCreate:           false,
			wantMarkTransactionCalled:   false,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.getByIdForUpdateError = domain.AccountNotFound
			},
		},

		{
			id:          1,
			name:        "Db error from transaction create",
			wantedError: ErrorDb,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantIncreaseBalance:         false,
			wantDepositCreate:           false,
			wantMarkTransactionCalled:   false,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				transactionRepo.createErr = ErrorDb
			},
		},

		{
			id:          1,
			name:        "Db error from increase balance",
			wantedError: ErrorDb,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantIncreaseBalance:         true,
			wantDepositCreate:           false,
			wantMarkTransactionCalled:   false,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.increaseBalanceError = ErrorDb
			},
		},

		{
			id:          1,
			name:        "Account not found error from increase balance",
			wantedError: domain.AccountNotFound,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantIncreaseBalance:         true,
			wantDepositCreate:           false,
			wantMarkTransactionCalled:   false,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				accountRepo.increaseBalanceError = domain.AccountNotFound
			},
		},

		{
			id:          1,
			name:        "Deposit create db error",
			wantedError: ErrorDb,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantIncreaseBalance:         true,
			wantDepositCreate:           true,
			wantMarkTransactionCalled:   false,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				depoRepo.createError = ErrorDb
			},
		},

		{
			id:          1,
			name:        "Mark transaction error transaction not found",
			wantedError: domain.ErrorTransactionNotFound,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantIncreaseBalance:         true,
			wantDepositCreate:           true,
			wantMarkTransactionCalled:   true,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				transactionRepo.markTransactionErr = domain.ErrorTransactionNotFound
			},
		},

		{
			id:          1,
			name:        "Mark transaction error db",
			wantedError: ErrorDb,
			req: dto.CreateDepositRequest{
				Amount: decimal.NewFromInt(10000),
				Source: "terminal",
			},
			wantGetByIdForUpdateCalled:  true,
			wantTransactionCreateCalled: true,
			wantIncreaseBalance:         true,
			wantDepositCreate:           true,
			wantMarkTransactionCalled:   true,

			setup: func(txRepo *fakeTxManager, depoRepo *fakeDepositRepo, accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo) {
				transactionRepo.markTransactionErr = ErrorDb
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			txManager := &fakeTxManager{}
			depoRepo := &fakeDepositRepo{}
			accountRepo := &fakeAccountRepo{
				account: domain.Account{
					Id:         test.id,
					IsActive:   true,
					Balance:    decimal.NewFromInt(10000),
					CurrencyId: 1,
					UserId:     12,
					CreatedAt:  time.Now(),
				},
			}
			transactionRepo := &fakeTransactionRepo{}

			test.setup(txManager, depoRepo, accountRepo, transactionRepo)

			srv := NewDepositService(depoRepo, accountRepo, txManager, transactionRepo)

			err := srv.Create(context.Background(), test.req, test.id)

			if !errors.Is(err, test.wantedError) {
				t.Fatalf("Expected error %v, got %v", test.wantedError, err)
			}

			if test.wantGetByIdForUpdateCalled != accountRepo.calledGetByIdForUpdate {
				t.Fatalf("Expected wantGetByIdForUpdateCalled: %v, got %v", test.wantGetByIdForUpdateCalled, accountRepo.calledGetByIdForUpdate)
			}

			if test.wantDepositCreate != depoRepo.createCalled {
				t.Fatalf("Expected wantDepositCreate: %v, got %v", test.wantDepositCreate, depoRepo.createCalled)
			}

			if test.wantIncreaseBalance != accountRepo.calledIncreaseBalance {
				t.Fatalf("Expected wantIncreaseBalance: %v, got %v", test.wantIncreaseBalance, accountRepo.calledIncreaseBalance)
			}

			if test.wantTransactionCreateCalled != transactionRepo.calledCreate {
				t.Fatalf("Expected wantTransactionCreateCalled: %v, got %v", test.wantTransactionCreateCalled, transactionRepo.calledCreate)
			}

			if test.wantMarkTransactionCalled != transactionRepo.calledMarkTransaction {
				t.Fatalf("Expected wantMarkTransactionCalled: %v, got %v", test.wantMarkTransactionCalled, transactionRepo.calledMarkTransaction)
			}
		})
	}

}
