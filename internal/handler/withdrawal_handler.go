package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/cache"
	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/middleware"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/alanzhumalin/bank/pkg/response"
	"github.com/rs/zerolog"
)

type withdrawalHandler struct {
	service               service.WithdrawalService
	idempotencyRedisStore cache.IdempotencyStoreInterface
	logger                zerolog.Logger
}

func NewWithDrawalHandler(idempotencyRedisStore cache.IdempotencyStoreInterface, service service.WithdrawalService, logger zerolog.Logger) *withdrawalHandler {
	return &withdrawalHandler{idempotencyRedisStore: idempotencyRedisStore, service: service, logger: logger.With().Str("component", "withdrawal_handler").Logger()}
}

func WithdrawalRouter(w *withdrawalHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /", middleware.Chain(http.HandlerFunc(w.Create), auth))

	return mux
}

func (wh *withdrawalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWindrawalRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, ok := r.Context().Value(dto.UserKey{}).(int)
	if !ok || userId <= 0 {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ok, res, err := wh.idempotencyRedisStore.Start(r.Context(), req.IdempotencyKey, dto.IdempotencyResponse{
		Status:   "pending",
		Response: []byte(`{}`),
	})
	if err != nil {
		wh.logger.Error().Err(err).Msg("Error in starting new key redis store")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !ok {
		wh.logger.Info().Str("account_id", strconv.Itoa(req.AccountId)).Str("amount", req.Amount.String()).Msg("Created withdraw")
		switch res.Status {
		case "pending":
			response.WriteError(w, http.StatusTooEarly, res)
		case "failed":
			response.WriteError(w, http.StatusConflict, res)
		case "completed":
			response.WriteJson(w, http.StatusOK, res)
		default:
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	idem, err := wh.service.Create(r.Context(), req, userId)
	if err != nil {

		switch err {
		case domain.ErrorIdempotencyPending:

			b := response.IdemErrorResponse(err)
			response.WriteError(w, http.StatusTooEarly, b)
			return
		case domain.ErrorIdempotencyFailed:
			b := response.IdemErrorResponse(err)

			response.WriteError(w, http.StatusConflict, b)
		case domain.AccountIsNotActive:
			b := response.IdemErrorResponse(err)
			response.WriteError(w, http.StatusConflict, b)
		case domain.ErrorNotEnoughBalance:
			b := response.IdemErrorResponse(err)

			response.WriteError(w, http.StatusConflict, b)

		case domain.AccountNotFound:
			b := response.IdemErrorResponse(err)

			response.WriteError(w, http.StatusNotFound, b)
		default:
			wh.logger.Error().Err(err).Msg("Error in create withdrawal")
			_ = wh.idempotencyRedisStore.Delete(r.Context(), req.IdempotencyKey)
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		responseByte, erro := json.Marshal(map[string]any{
			"data": err.Error(),
		})

		if erro != nil {
			wh.logger.Error().Err(err).Msg("Error in creating json")
			return
		}

		if err := wh.idempotencyRedisStore.Failed(r.Context(), req.IdempotencyKey, dto.IdempotencyResponse{
			Status:   "failed",
			Response: responseByte,
		}); err != nil {
			wh.logger.Error().Err(err).Msg("Error in failed redis key")
		}
		return

	}

	if err = wh.idempotencyRedisStore.Complete(r.Context(), req.IdempotencyKey, idem); err != nil {
		wh.logger.Error().Err(err).Msg("Error in completing redis key")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	wh.logger.Info().Str("account_id", strconv.Itoa(req.AccountId)).Str("amount", req.Amount.String()).Msg("Created withdraw")
	response.WriteJson(w, http.StatusCreated, idem)

}
