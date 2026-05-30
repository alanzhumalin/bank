package service

import (
	"context"
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
		return true, cr.existsErr
	}

	return false, nil
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
