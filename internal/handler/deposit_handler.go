package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/rs/zerolog"
)

type depositHandler struct {
	service service.DepositService

	logger zerolog.Logger
}

func NewDepositHandler(s service.DepositService, l zerolog.Logger) *depositHandler {
	return &depositHandler{service: s, logger: l.With().Str("component", "deposit_handler").Logger()}
}

func DepositRouter(depositHandler *depositHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /{account_id}", depositHandler.Create)

	return mux
}

func (d *depositHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDepositRequest

	pathId := r.PathValue("account_id")

	if pathId == "" {
		WriteError(w, http.StatusBadRequest, "account id is required")
		return
	}

	id, err := strconv.Atoi(pathId)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = d.service.Create(r.Context(), req, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.AccountNotFound):
			WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.AccountIsNotActive):
			WriteError(w, http.StatusNotFound, err.Error())
		default:
			d.logger.Error().Err(err).Str("account_id", pathId).Msg("Error occured")
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	d.logger.Info().Str("account_id", pathId).Str("amount", req.Amount.String()).Msg("Created deposit")
	WriteJson(w, http.StatusCreated, "created")
}
