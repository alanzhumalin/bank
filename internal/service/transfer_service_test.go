package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/pkg/pagination"
	"github.com/shopspring/decimal"
)

// type transferService struct {
// 	transferRepo    repository.TransferRepository
// 	txManager       repository.TxManagerRepository
// 	accountRepo     repository.AccountRepository
// 	transactionRepo repository.TransactionRepository
// }

var (
	ErrorDb = errors.New("Error db")
)

type fakeTransferRepo struct {
	createErr     error
	createdCalled bool
}

func (tR *fakeTransferRepo) Create(ctx context.Context, t ...domain.Transfer) error {
	tR.createdCalled = true
	if tR.createErr != nil {
		return tR.createErr
	}
	return nil
}

type fakeTxManager struct {
	withTxError error
	calledTx    bool
}

func (tx *fakeTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx.calledTx = true

	if tx.withTxError != nil {
		return tx.withTxError
	}

	return fn(ctx)
}

type fakeAccountRepo struct {
	account1 domain.Account
	account2 domain.Account

	account domain.Account

	selectTwoAccountsForUpdateError error
	increaseBalanceError            error
	decreaseBalanceError            error
	getByIdForUpdateError           error

	calledSelectTwoAccountForUpdate bool
	calledIncreaseBalance           bool
	calledDecreaseBalance           bool
	calledGetByIdForUpdate          bool
}

func (ac *fakeAccountRepo) SelectTwoAccountsForUpdate(ctx context.Context, senderAccountId int, receiverAccountId int) (domain.Account, domain.Account, error) {
	ac.calledSelectTwoAccountForUpdate = true

	if ac.selectTwoAccountsForUpdateError != nil {
		return domain.Account{}, domain.Account{}, ac.selectTwoAccountsForUpdateError
	}
	return ac.account1, ac.account2, nil
}

func (ac *fakeAccountRepo) IncreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error {
	ac.calledIncreaseBalance = true

	if ac.increaseBalanceError != nil {
		return ac.increaseBalanceError
	}

	return nil
}

func (ac *fakeAccountRepo) DecreaseBalance(ctx context.Context, balance decimal.Decimal, accountId int) error {
	ac.calledDecreaseBalance = true

	if ac.decreaseBalanceError != nil {
		return ac.decreaseBalanceError
	}

	return nil
}

func (ac *fakeAccountRepo) DeleteById(ctx context.Context, id int, time time.Time) error {
	return nil
}

func (ac *fakeAccountRepo) GetAll(ctx context.Context) ([]domain.Account, error) {
	return []domain.Account{}, nil
}

func (ac *fakeAccountRepo) GetByIdForUpdate(ctx context.Context, id int) (domain.Account, error) {
	ac.calledGetByIdForUpdate = true

	if ac.getByIdForUpdateError != nil {
		return domain.Account{}, ac.getByIdForUpdateError
	}
	return ac.account, nil
}

func (ac *fakeAccountRepo) GetUserAccounts(ctx context.Context, userId int) ([]domain.Account, error) {
	return []domain.Account{}, nil

}

func (ac *fakeAccountRepo) Exists(ctx context.Context, userId int, currencyId int) (bool, error) {
	return true, nil
}

func (ac *fakeAccountRepo) Create(ctx context.Context, a domain.Account) error {
	return nil
}

type fakeTransactionRepo struct {
	markTransactionErr error
	createErr          error

	calledCreate          bool
	calledMarkTransaction bool
}

func (f *fakeTransactionRepo) Create(ctx context.Context, t ...domain.Transaction) (map[int]int, error) {
	f.calledCreate = true

	if f.createErr != nil {
		return map[int]int{}, f.createErr
	}

	mp := make(map[int]int, len(t))

	for i, v := range t {
		mp[v.AccountId] = 100 + i
	}

	return mp, nil
}

func (f *fakeTransactionRepo) MarkTransaction(ctx context.Context, status string, status_message string, id int) error {
	f.calledMarkTransaction = true

	if f.markTransactionErr != nil {
		return f.markTransactionErr
	}

	return nil
}
func (f *fakeTransactionRepo) GetByAccountId(ctx context.Context, id int, limit int, transactionCursor *pagination.TransactionCursor) ([]domain.Transaction, int, error) {
	return []domain.Transaction{}, 0, nil
}

func (f *fakeTransactionRepo) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func (f *fakeTransactionRepo) GetByUserId(ctx context.Context, userId int, cursor *pagination.TransactionCursor, limit int, currencies *[]string) ([]domain.Transaction, error) {
	return []domain.Transaction{}, nil
}

func TestTransferService_CreateSuccess(t *testing.T) {
	fakeTransferRepo := &fakeTransferRepo{}
	fakeTxManager := &fakeTxManager{}
	fakeAccountRepo := &fakeAccountRepo{
		account1: domain.Account{
			Id:         100,
			UserId:     1,
			CurrencyId: 1,
			Balance:    decimal.NewFromInt(39999),
			IsActive:   true,
			CreatedAt:  time.Now(),
		},
		account2: domain.Account{
			Id:         101,
			UserId:     2,
			CurrencyId: 1,
			Balance:    decimal.NewFromInt(39999),
			IsActive:   true,
			CreatedAt:  time.Now(),
		},
	}
	fakeTransactionRepo := &fakeTransactionRepo{}

	srv := NewTransferService(fakeTransferRepo, fakeTxManager, fakeAccountRepo, fakeTransactionRepo)

	err := srv.Create(context.Background(), dto.CreateTransferRequest{
		SenderAccountId:   100,
		ReceiverAccountId: 101,
		CurrencyId:        1,
		Amount:            decimal.NewFromInt(10000),
	})

	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !fakeTxManager.calledTx {
		t.Error("Expected txmanager to be called")
	}

	if !fakeAccountRepo.calledSelectTwoAccountForUpdate {
		t.Error("Expected selectwoaccountforupdate to be called")
	}

	if !fakeTransactionRepo.calledCreate {
		t.Error("Expected create transaction to be called")
	}

	if !fakeAccountRepo.calledDecreaseBalance {
		t.Error("Expected decreaseBalance transaction to be called")
	}

	if !fakeAccountRepo.calledIncreaseBalance {
		t.Error("Expected increaseBalance transaction to be called")
	}

	if !fakeTransferRepo.createdCalled {
		t.Error("Expected transfer create to be called")
	}

	if !fakeTransactionRepo.calledMarkTransaction {
		t.Error("Expected transaction marked to be called")
	}

}

func TestTransferErrors(t *testing.T) {
	tests := []struct {
		sender      domain.Account
		receiver    domain.Account
		req         dto.CreateTransferRequest
		name        string //testname
		wantedError error
	}{
		{
			sender: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},
			receiver: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 2,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(1000),
			},
			name:        "receiver currency mismatch",
			wantedError: domain.AccountNotSupportCurrency,
		}, {
			sender: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},
			receiver: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 2,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   false,
				CreatedAt:  time.Now(),
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        2,
				Amount:            decimal.NewFromInt(1000),
			},
			name:        "receiver account is not active",
			wantedError: domain.AccountIsNotActive,
		},

		{
			sender: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},
			receiver: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 2,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        2,
				Amount:            decimal.NewFromInt(1000),
			},
			name:        "sender currency mismatch",
			wantedError: domain.AccountNotSupportCurrency,
		},

		{
			sender: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   false,
				CreatedAt:  time.Now(),
			},
			receiver: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},

			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(1000),
			},
			wantedError: domain.AccountIsNotActive,
			name:        "sender account is not active",
		},

		{
			sender: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(999),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},
			receiver: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
				CreatedAt:  time.Now(),
			},

			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(1000),
			},
			wantedError: domain.ErrorNotEnoughBalance,
			name:        "sender has less money",
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			t.Parallel()
			transferRepo := &fakeTransferRepo{}
			txManager := &fakeTxManager{}
			transactionRepo := &fakeTransactionRepo{}
			accountRepo := &fakeAccountRepo{
				account1: test.sender,
				account2: test.receiver,
			}

			srv := transferService{
				transferRepo:    transferRepo,
				txManager:       txManager,
				accountRepo:     accountRepo,
				transactionRepo: transactionRepo,
			}

			err := srv.Create(context.Background(), test.req)

			if !errors.Is(err, test.wantedError) {
				t.Errorf("Expected error %v, got %v", test.wantedError, err)
			}

		})

	}
}

func TestTransferServiceRepoErrors(t *testing.T) {
	tests := []struct {
		name                                 string
		wantedError                          error
		account1                             domain.Account
		account2                             domain.Account
		req                                  dto.CreateTransferRequest
		wantTransactionCreateCalled          bool
		wantDecreaseBalanceCalled            bool
		wantIncreaseBalanceCalled            bool
		wantTransferCreateCalled             bool
		wantMarkTransactionCalled            bool
		wantSelectTwoAccountsForUpdateCalled bool
		setup                                func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo)
	}{
		{
			name:                                 "select two accounts error",
			wantedError:                          domain.AccountNotFound,
			wantTransactionCreateCalled:          false,
			wantDecreaseBalanceCalled:            false,
			wantIncreaseBalanceCalled:            false,
			wantTransferCreateCalled:             false,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				accountRepo.selectTwoAccountsForUpdateError = domain.AccountNotFound
			},
		},

		{
			name:                                 "create transaction error",
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            false,
			wantIncreaseBalanceCalled:            false,
			wantTransferCreateCalled:             false,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			wantedError:                          ErrorDb,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				transactionRepo.createErr = ErrorDb
			},
		},

		{
			name:                                 "decrease balance, account not found error",
			wantedError:                          domain.AccountNotFound,
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            false,
			wantTransferCreateCalled:             false,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				accountRepo.decreaseBalanceError = domain.AccountNotFound
			},
		},

		{
			name:                                 "decrease balance, error db",
			wantedError:                          ErrorDb,
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            false,
			wantTransferCreateCalled:             false,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				accountRepo.decreaseBalanceError = ErrorDb
			},
		},

		{
			name:                                 "increase balance, account not found error",
			wantedError:                          domain.AccountNotFound,
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            true,
			wantTransferCreateCalled:             false,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				accountRepo.increaseBalanceError = domain.AccountNotFound
			},
		},

		{
			name:                                 "increase balance, error db",
			wantedError:                          ErrorDb,
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            true,
			wantTransferCreateCalled:             false,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				accountRepo.increaseBalanceError = ErrorDb
			},
		},

		{
			name:                                 "create transfer error db",
			wantedError:                          ErrorDb,
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            true,
			wantTransferCreateCalled:             true,
			wantMarkTransactionCalled:            false,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				transferRepo.createErr = ErrorDb
			},
		},

		{
			name:        "mark transaction, transaction not found error",
			wantedError: domain.ErrorTransactionNotFound,

			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            true,
			wantTransferCreateCalled:             true,
			wantMarkTransactionCalled:            true,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				transactionRepo.markTransactionErr = domain.ErrorTransactionNotFound
			},
		},

		{
			name:                                 "mark transaction, db error",
			wantedError:                          ErrorDb,
			wantTransactionCreateCalled:          true,
			wantDecreaseBalanceCalled:            true,
			wantIncreaseBalanceCalled:            true,
			wantTransferCreateCalled:             true,
			wantMarkTransactionCalled:            true,
			wantSelectTwoAccountsForUpdateCalled: true,
			account1: domain.Account{
				Id:         100,
				UserId:     1,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			account2: domain.Account{
				Id:         101,
				UserId:     2,
				CurrencyId: 1,
				Balance:    decimal.NewFromInt(10000),
				IsActive:   true,
			},
			req: dto.CreateTransferRequest{
				SenderAccountId:   100,
				ReceiverAccountId: 101,
				CurrencyId:        1,
				Amount:            decimal.NewFromInt(100),
			},
			setup: func(accountRepo *fakeAccountRepo, transactionRepo *fakeTransactionRepo, txRepo *fakeTxManager, transferRepo *fakeTransferRepo) {
				transactionRepo.markTransactionErr = ErrorDb
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			transferRepo := &fakeTransferRepo{}
			txManager := &fakeTxManager{}
			accountRepo := &fakeAccountRepo{
				account1: test.account1,
				account2: test.account2,
			}
			transactionRepo := &fakeTransactionRepo{}

			test.setup(accountRepo, transactionRepo, txManager, transferRepo)

			srv := transferService{
				transferRepo:    transferRepo,
				txManager:       txManager,
				accountRepo:     accountRepo,
				transactionRepo: transactionRepo,
			}

			err := srv.Create(context.Background(), test.req)

			if !errors.Is(err, test.wantedError) {
				t.Fatalf("Expected error %v, got %v", test.wantedError, err)
			}

			// wantTransactionCreateCalled          bool
			// wantDecreaseBalanceCalled            bool
			// wantIncreaseBalanceCalled            bool
			// wantTransferCreateCalled             bool
			// wantMarkTransactionCalled            bool
			// wantSelectTwoAccountsForUpdateCalled bool

			if transactionRepo.calledCreate != test.wantTransactionCreateCalled {
				t.Fatalf("Expected calledCreate %v, got %v", transactionRepo.calledCreate, test.wantTransactionCreateCalled)
			}

			if accountRepo.calledDecreaseBalance != test.wantDecreaseBalanceCalled {
				t.Fatalf("Expected calledDecreaseBalance %v, got %v", accountRepo.calledDecreaseBalance, test.wantDecreaseBalanceCalled)
			}

			if accountRepo.calledIncreaseBalance != test.wantIncreaseBalanceCalled {
				t.Fatalf("Expected calledIncreaseBalance %v, got %v", accountRepo.calledIncreaseBalance, test.wantIncreaseBalanceCalled)
			}

			if transferRepo.createdCalled != test.wantTransferCreateCalled {
				t.Fatalf("Expected transferCreate %v, got %v", transferRepo.createdCalled, test.wantTransferCreateCalled)
			}

			if transactionRepo.calledMarkTransaction != test.wantMarkTransactionCalled {
				t.Fatalf("Expected calledMarkTransaction %v, got %v", transactionRepo.calledMarkTransaction, test.wantMarkTransactionCalled)
			}

			if accountRepo.calledSelectTwoAccountForUpdate != test.wantSelectTwoAccountsForUpdateCalled {
				t.Fatalf("Expected calledSelectTwoAccountForUpdate %v, got %v", accountRepo.calledSelectTwoAccountForUpdate, test.wantSelectTwoAccountsForUpdateCalled)
			}

		})
	}
}
