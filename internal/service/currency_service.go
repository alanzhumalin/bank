package service

import (
	"context"
	"fmt"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/repository"
)

type currencyService struct {
	repo repository.CurrencyRepository
}

func NewCurrencyService(repo repository.CurrencyRepository) CurrencyService {
	return &currencyService{repo: repo}
}

func (c *currencyService) Create(ctx context.Context, req dto.CreateNewCurrencyRequest) error {
	check, err := c.repo.Exists(ctx, req.Code)

	if err != nil {
		return fmt.Errorf("Error create currency in currency_service: %w", err)
	}

	if check {
		return domain.ErrorCurrencyAlreadyExists
	}

	currency := domain.NewCurrency(req.Name, req.Code, req.Symbol)

	err = c.repo.Create(ctx, currency)

	if err != nil {
		return fmt.Errorf("Error create in currency_service: %w", err)
	}

	return nil
}

func (c *currencyService) Delete(ctx context.Context, id int) error {
	err := c.repo.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf("Error delete currency in currency_service: %w", err)
	}

	return nil
}

func (c *currencyService) Update(ctx context.Context, id int, req dto.UpdateCurrency) error {
	err := c.repo.UpdateById(ctx, id, req.Name, req.Code, req.Symbol)

	if err != nil {
		return fmt.Errorf("Error update currency in currency_service: %w", err)
	}

	return nil
}

func (c *currencyService) GetAll(ctx context.Context) ([]dto.GetCurrencyResponse, error) {
	res, err := c.repo.GetAll(ctx)

	if err != nil {
		return []dto.GetCurrencyResponse{}, fmt.Errorf("Error get all in currency_service: %w", err)
	}

	currencies := make([]dto.GetCurrencyResponse, 0, len(res))

	for _, cur := range res {
		currencies = append(currencies, dto.NewGetCurrencyResponse(cur))
	}

	return currencies, nil
}

func (c *currencyService) GetById(ctx context.Context, id int) (dto.GetCurrencyResponse, error) {

	cur, err := c.repo.GetById(ctx, id)

	if err != nil {
		return dto.GetCurrencyResponse{}, fmt.Errorf("Error get by id in currency_service: %w", err)
	}

	r := dto.NewGetCurrencyResponse(cur)

	return r, nil
}
