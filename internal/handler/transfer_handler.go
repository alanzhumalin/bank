package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/rs/zerolog"
)

type transferHandler struct {
	service service.TransferService
	logger  zerolog.Logger
}

func NewTransferHandler(s service.TransferService, l zerolog.Logger) *transferHandler {
	return &transferHandler{service: s, logger: l.With().Str("component", "transfer_handler").Logger()}
}

func TransferRouter(t *transferHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", t.Create)

	return mux
}

func (t *transferHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransferRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		WriteJson(w, http.StatusBadRequest, err.Error())
	}

	err = t.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorNotEnoughBalance):
			WriteError(w, http.StatusOK, err.Error())
		default:
			t.logger.Error().Err(err).Msg("Error in creating transfer")
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}
	t.logger.Info().Msg("Created transfer")
	WriteJson(w, http.StatusCreated, "created")
}
