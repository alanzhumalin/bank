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

type fakeIdempotencyRedis struct {
}

func (f *fakeIdempotencyRedis) Start(ctx context.Context, key string, res dto.IdempotencyResponse) (bool, dto.IdempotencyResponse, error) {
	return true, dto.IdempotencyResponse{}, nil
}
func (f *fakeIdempotencyRedis) Complete(ctx context.Context, key string, res dto.IdempotencyResponse) error {
	return nil
}

func (f *fakeIdempotencyRedis) Failed(ctx context.Context, key string, res dto.IdempotencyResponse) error {
	return nil
}
func (f *fakeIdempotencyRedis) Delete(ctx context.Context, key string) error {
	return nil
}

type fakeWithdrawalService struct {
	createCalled bool
	createError  error
}

var (
	DbError = errors.New("Db error")
)

func (w *fakeWithdrawalService) Create(ctx context.Context, req dto.CreateWindrawalRequest, userId int) (dto.IdempotencyResponse, error) {
	w.createCalled = true

	if w.createError != nil {
		return dto.IdempotencyResponse{}, w.createError
	}

	return dto.IdempotencyResponse{}, nil
}

func TestWithdrawalHandlerCreateSuccess(t *testing.T) {
	withdrawalService := &fakeWithdrawalService{}
	idempotencyStore := &fakeIdempotencyRedis{}
	logger := zerolog.Nop()

	handler := NewWithDrawalHandler(idempotencyStore, withdrawalService, logger)

	body := strings.NewReader(`
		{
			"account_id": 1,
			"amount": 123123, 
			"source": "terminal",
			"idempotency_key": "msodfmomomo1m3o21m3"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), dto.UserKey{}, 31)
	req = req.WithContext(ctx)

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
	idempotencyStore := &fakeIdempotencyRedis{}
	logger := zerolog.Nop()

	body := strings.NewReader(`{bad json`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()

	handler := NewWithDrawalHandler(idempotencyStore, withdrawalService, logger)

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
	idempotencyStore := &fakeIdempotencyRedis{}

	logger := zerolog.Nop()

	body := strings.NewReader(`
		{
			"amount": 123123, 
			"source": "terminal",
			"idempotency_key": "msodfmomomo1m3o21m3"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()

	handler := NewWithDrawalHandler(idempotencyStore, withdrawalService, logger)

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
	idempotencyStore := &fakeIdempotencyRedis{}

	logger := zerolog.Nop()

	body := strings.NewReader(`
		{
			"account_id": 1,
			"amount": 123123, 
			"source": "terminal",
			"idempotency_key": "msodfmomomo1m3o21m3"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/withdrawals/", body)
	res := httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), dto.UserKey{}, 31)
	req = req.WithContext(ctx)

	handler := NewWithDrawalHandler(idempotencyStore, withdrawalService, logger)

	handler.Create(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("Expected code %v, got %v, body: %v", http.StatusBadRequest, res.Code, res.Body.String())
	}

	if !withdrawalService.createCalled {
		t.Fatalf("create function must be called")
	}

}
