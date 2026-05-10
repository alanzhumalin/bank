package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/middleware"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/alanzhumalin/bank/pkg/response"
	"github.com/rs/zerolog"
)

type withdrawalHandler struct {
	service service.WithdrawalService
	logger  zerolog.Logger
}

func NewWithDrawalHandler(service service.WithdrawalService, logger zerolog.Logger) *withdrawalHandler {
	return &withdrawalHandler{service: service, logger: logger.With().Str("component", "withdrawal_handler").Logger()}
}

func WithdrawalRouter(w *withdrawalHandler, authMiddleware middleware.AuthMiddleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /", middleware.Chain(http.HandlerFunc(w.Create), authMiddleware.Middleware()))

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

	if err = wh.service.Create(r.Context(), req); err != nil {
		wh.logger.Error().Err(err).Msg("Error in create withdrawal")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	wh.logger.Info().Str("account_id", strconv.Itoa(req.AccountId)).Str("amount", req.Amount.String()).Msg("Created withdraw")
	response.WriteJson(w, http.StatusCreated, "created")

}
