package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/rs/zerolog"
)

type fakeWithdrawalService struct {
	createCalled bool
	createError  error
}

var (
	DbError = errors.New("Db error")
)

func (w *fakeWithdrawalService) Create(ctx context.Context, req dto.CreateWindrawalRequest) error {
	w.createCalled = true

	if w.createError != nil {
		return w.createError
	}

	return nil
}

func TestWithdrawalHandlerCreateSuccess(t *testing.T) {
	withdrawalService := &fakeWithdrawalService{}
	logger := zerolog.Nop()

	handler := NewWithDrawalHandler(withdrawalService, logger)

	body := strings.NewReader(`
		{
			"account_id": 1,
			"amount": 123123, 
			"source": "terminal"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()

	handler.Create(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusCreated, res.Code, res.Body.String())
	}

	if !withdrawalService.createCalled {
		t.Fatal("Expected create withdrawal to be called")
	}
}

func TestWithdrawalCreateInvalidJson(t *testing.T) {
	withdrawalService := &fakeWithdrawalService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{bad json`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()

	handler := NewWithDrawalHandler(withdrawalService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if withdrawalService.createCalled {
		t.Fatalf("create function must not be called")
	}

}

func TestWithdrawalCreateMissingAccountId(t *testing.T) {
	withdrawalService := &fakeWithdrawalService{}
	logger := zerolog.Nop()

	body := strings.NewReader(`
		{
			"amount": 123123, 
			"source": "terminal"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()

	handler := NewWithDrawalHandler(withdrawalService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if withdrawalService.createCalled {
		t.Fatalf("create function must not be called")
	}

}

func TestWithdrawalCreateError(t *testing.T) {
	withdrawalService := &fakeWithdrawalService{
		createError: DbError,
	}
	logger := zerolog.Nop()

	body := strings.NewReader(`
		{
			"account_id": 1,
			"amount": 123123, 
			"source": "terminal"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()

	handler := NewWithDrawalHandler(withdrawalService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if !withdrawalService.createCalled {
		t.Fatalf("create function must be called")
	}

}
