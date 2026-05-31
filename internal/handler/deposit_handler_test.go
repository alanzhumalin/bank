package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/rs/zerolog"
)

type fakeDepositService struct {
	createCalled bool
	createErr    error
}

func (f *fakeDepositService) Create(ctx context.Context, req dto.CreateDepositRequest, id int) error {
	f.createCalled = true

	if f.createErr != nil {
		return f.createErr
	}

	return nil
}

func TestDepositHandlerCreateSuccess(t *testing.T) {
	depositService := &fakeDepositService{}
	logger := zerolog.Nop()

	handler := NewDepositHandler(depositService, logger)
	body := strings.NewReader(`{
		"amount": 14000,
		"source": "terminal"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/deposit/1", body)
	req.SetPathValue("account_id", "1")

	res := httptest.NewRecorder()

	handler.Create(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("Expected %v, got %v, body: %v", http.StatusCreated, res.Code, res.Body.String())
	}

	if !depositService.createCalled {
		t.Fatal("expected depositService.Create to be called")
	}

}

func TestDepositHandlerCreateInvalidAccountId(t *testing.T) {
	depositService := &fakeDepositService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{
		"amount": 14000,
		"source": "terminal"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/deposit/abc", body)
	req.SetPathValue("account_id", "abc")

	res := httptest.NewRecorder()

	handler := NewDepositHandler(depositService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if depositService.createCalled {
		t.Fatalf("The create function must not be called")
	}
}

func TestDepositHandlerCreateInvalidJson(t *testing.T) {
	depositService := &fakeDepositService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{bad json`)

	req := httptest.NewRequest(http.MethodPost, "/deposit/1", body)
	req.SetPathValue("account_id", "1")

	res := httptest.NewRecorder()

	handler := NewDepositHandler(depositService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if depositService.createCalled {
		t.Fatalf("The create function must not be called")
	}
}

func TestDepositHandlerCreateMissingAmount(t *testing.T) {
	depositService := &fakeDepositService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{
		"source": "terminal"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/deposit/1", body)
	req.SetPathValue("account_id", "1")
	res := httptest.NewRecorder()

	handler := NewDepositHandler(depositService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body)
	}

	if depositService.createCalled {
		t.Fatalf("The create function must not be called")
	}
}

func TestDepositHandlerCreateServiceError(t *testing.T) {
	depositService := &fakeDepositService{
		createErr: domain.AccountNotFound,
	}
	logger := zerolog.Nop()

	body := strings.NewReader(`{
		"amount": 123123,
		"source": "terminal"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/deposit/1", body)
	req.SetPathValue("account_id", "1")
	res := httptest.NewRecorder()

	handler := NewDepositHandler(depositService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body)
	}

	if !depositService.createCalled {
		t.Fatalf("The create function must be called")
	}
}
