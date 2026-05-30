package service

import (
	"context"
	"errors"
	"testing"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
)

type fakeCurrencyRepo struct {
	currency domain.Сurrency

	createdCalled    bool
	deleteCalled     bool
	getByIdCalled    bool
	updateByIdCalled bool
	existsCalled     bool
	getAllCalled     bool

	createErr     error
	deleteErr     error
	getByIdErr    error
	updateByIdErr error
	existsErr     error
	getAllErr     error

	exists bool
}

func (cr *fakeCurrencyRepo) Create(ctx context.Context, c domain.Сurrency) error {
	cr.createdCalled = true

	if cr.createErr != nil {
		return cr.createErr
	}

	return nil
}
func (cr *fakeCurrencyRepo) Delete(ctx context.Context, id int) error {
	cr.deleteCalled = true

	if cr.deleteErr != nil {
		return cr.deleteErr
	}

	return nil
}
func (cr *fakeCurrencyRepo) GetById(ctx context.Context, id int) (domain.Сurrency, error) {
	cr.getByIdCalled = true

	if cr.getByIdErr != nil {
		return domain.Сurrency{}, cr.getByIdErr
	}

	return cr.currency, nil
}
func (cr *fakeCurrencyRepo) UpdateById(ctx context.Context, id int, name string, code string, symbol string) error {
	cr.updateByIdCalled = true

	if cr.updateByIdErr != nil {
		return cr.updateByIdErr
	}

	return nil
}
func (cr *fakeCurrencyRepo) Exists(ctx context.Context, code string) (bool, error) {
	cr.existsCalled = true

	if cr.existsErr != nil {
		return false, cr.existsErr
	}

	return cr.exists, nil
}
func (cr *fakeCurrencyRepo) GetAll(ctx context.Context) ([]domain.Сurrency, error) {
	cr.getAllCalled = true

	if cr.getAllErr != nil {
		return []domain.Сurrency{}, cr.getAllErr
	}

	return []domain.Сurrency{
		cr.currency,
	}, nil
}

func TestCurrencyServiceCreateSuccess(t *testing.T) {
	currencyRepository := &fakeCurrencyRepo{}

	srv := NewCurrencyService(currencyRepository)

	req := dto.CreateNewCurrencyRequest{
		Name:   "American dollar",
		Code:   "USD",
		Symbol: "$",
	}

	err := srv.Create(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}

	if !currencyRepository.existsCalled {
		t.Fatal("Expected called exists")
	}

	if !currencyRepository.createdCalled {
		t.Fatal("Expected called create")
	}
}

func TestCurrencyServiceCreateErrors(t *testing.T) {
	tests := []struct {
		name         string
		wantedErr    error
		req          dto.CreateNewCurrencyRequest
		existsCalled bool
		createCalled bool
		setup        func(repo *fakeCurrencyRepo)
	}{
		{
			name:      "check for existence of currency",
			wantedErr: domain.ErrorCurrencyAlreadyExists,
			req: dto.CreateNewCurrencyRequest{
				Name:   "American dollar",
				Code:   "USD",
				Symbol: "$",
			},
			existsCalled: true,
			createCalled: false,
			setup: func(repo *fakeCurrencyRepo) {
				repo.exists = true
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			currencyRepo := &fakeCurrencyRepo{}

			test.setup(currencyRepo)

			srv := NewCurrencyService(currencyRepo)

			err := srv.Create(context.Background(), test.req)

			if !errors.Is(err, test.wantedErr) {
				t.Fatalf("Expected %v, got %v", test.wantedErr, err)
			}

			if test.createCalled != currencyRepo.createdCalled {
				t.Fatalf("Expected createdCalled %v, got %v", test.createCalled, currencyRepo.createdCalled)
			}

			if test.existsCalled != currencyRepo.existsCalled {
				t.Fatalf("Expected exists called %v, got %v", test.existsCalled, currencyRepo.existsCalled)
			}
		})
	}
}

func TestCurrencyServiceCreateRepoErrors(t *testing.T) {
	tests := []struct {
		name         string
		wantedErr    error
		req          dto.CreateNewCurrencyRequest
		calledExists bool
		calledCreate bool
		setup        func(repo *fakeCurrencyRepo)
	}{
		{
			name:      "Error in exists db",
			wantedErr: ErrorDb,
			req: dto.CreateNewCurrencyRequest{
				Name:   "American dollar",
				Code:   "USD",
				Symbol: "$",
			},
			calledExists: true,
			calledCreate: false,
			setup: func(repo *fakeCurrencyRepo) {
				repo.existsErr = ErrorDb
			},
		},

		{
			name:      "Error in create currency db",
			wantedErr: ErrorDb,
			req: dto.CreateNewCurrencyRequest{
				Name:   "American dollar",
				Code:   "USD",
				Symbol: "$",
			},
			calledExists: true,
			calledCreate: true,
			setup: func(repo *fakeCurrencyRepo) {
				repo.createErr = ErrorDb
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			currencyRepository := &fakeCurrencyRepo{}

			test.setup(currencyRepository)

			srv := NewCurrencyService(currencyRepository)

			err := srv.Create(context.Background(), test.req)

			if !errors.Is(err, test.wantedErr) {
				t.Fatalf("Expected %v, got %v", test.wantedErr, err)
			}

			if test.calledExists != currencyRepository.existsCalled {
				t.Fatalf("Expected calledExists %v, got %v", test.calledExists, currencyRepository.existsCalled)
			}

			if test.calledCreate != currencyRepository.createdCalled {
				t.Fatalf("Expected calledCreate %v, got %v", test.calledCreate, currencyRepository.createdCalled)
			}

		})
	}
}
