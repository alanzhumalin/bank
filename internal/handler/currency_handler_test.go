package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/rs/zerolog"
)

type fakeCurrencyService struct {
	createCalled  bool
	deleteCalled  bool
	updateCalled  bool
	getAllCalled  bool
	getByIdCalled bool

	createError  error
	deleteError  error
	updateError  error
	getAllError  error
	getByIdError error
}

var (
	DbError = errors.New("Db error")
)

func (c *fakeCurrencyService) Create(ctx context.Context, req dto.CreateNewCurrencyRequest) error {
	c.createCalled = true

	if c.createError != nil {
		return c.createError
	}
	return nil
}
func (c *fakeCurrencyService) Delete(ctx context.Context, id int) error {
	c.deleteCalled = true

	if c.deleteError != nil {
		return c.deleteError
	}
	return nil
}
func (c *fakeCurrencyService) Update(ctx context.Context, id int, req dto.UpdateCurrency) error {
	c.updateCalled = true

	if c.updateError != nil {
		return c.updateError
	}
	return nil
}
func (c *fakeCurrencyService) GetAll(ctx context.Context) ([]dto.GetCurrencyResponse, error) {
	c.getAllCalled = true

	if c.getAllError != nil {
		return []dto.GetCurrencyResponse{}, c.getAllError
	}
	return []dto.GetCurrencyResponse{}, nil
}
func (c *fakeCurrencyService) GetById(ctx context.Context, id int) (dto.GetCurrencyResponse, error) {
	c.getByIdCalled = true

	if c.getByIdError != nil {
		return dto.GetCurrencyResponse{}, c.getByIdError
	}
	return dto.GetCurrencyResponse{}, nil
}

func TestCurrencyHandlerCreateSuccess(t *testing.T) {
	currencyService := &fakeCurrencyService{}
	logger := zerolog.Nop()
	body := strings.NewReader(`{
		"name": "Dollar",
		"code": "USD",
		"symbol": "$"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/currencies/", body)
	res := httptest.NewRecorder()

	handler := NewCurrencyHandler(currencyService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusCreated, res.Code, res.Body.String())
	}

	if !currencyService.createCalled {
		t.Fatal("create function must be called")
	}
}

func TestCurrencyHandlerCreateInvalidJson(t *testing.T) {
	fakeCurrencyService := &fakeCurrencyService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{bad json`)

	req := httptest.NewRequest(http.MethodPost, "/currencies/", body)
	res := httptest.NewRecorder()

	handler := NewCurrencyHandler(fakeCurrencyService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if fakeCurrencyService.createCalled {
		t.Fatal("Create function must not be called")
	}
}

func TestCurrencyHandlerCreateMissingName(t *testing.T) {
	fakeCurrencyService := &fakeCurrencyService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{
		"code": "USD",
		"symbol": "$"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/currencies/", body)
	res := httptest.NewRecorder()

	handler := NewCurrencyHandler(fakeCurrencyService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if fakeCurrencyService.createCalled {
		t.Fatal("Create function must not be called")
	}
}

func TestCurrencyHandlerCreateError(t *testing.T) {

	tests := []struct {
		name             string
		wantedErr        error
		wantedStatusCode int
		request          *strings.Reader
		setup            func(s *fakeCurrencyService)
	}{
		{
			name:             "Currency already exists error",
			wantedErr:        domain.ErrorCurrencyAlreadyExists,
			wantedStatusCode: http.StatusConflict,
			request: strings.NewReader(`{
				"name": "Dollar",
				"code": "USD",
				"symbol": "$"
			}`),
			setup: func(s *fakeCurrencyService) {
				s.createError = domain.ErrorCurrencyAlreadyExists
			},
		},

		{
			name:             "Db error",
			wantedErr:        DbError,
			wantedStatusCode: http.StatusInternalServerError,
			request: strings.NewReader(`{
				"name": "Dollar",
				"code": "USD",
				"symbol": "$"
			}`),
			setup: func(s *fakeCurrencyService) {
				s.createError = DbError
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {

			fakeCurrencyService := &fakeCurrencyService{}

			test.setup(fakeCurrencyService)

			logger := zerolog.Nop()

			req := httptest.NewRequest(http.MethodPost, "/currencies/", test.request)
			res := httptest.NewRecorder()

			handler := NewCurrencyHandler(fakeCurrencyService, logger)

			handler.Create(res, req)

			if res.Code != test.wantedStatusCode {
				t.Fatalf("Expected status %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
			}

			if !fakeCurrencyService.createCalled {
				t.Fatal("Create function must be called")
			}
		})
	}

}
