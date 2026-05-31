package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/alanzhumalin/bank/pkg/response"
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
		response.WriteError(w, http.StatusBadRequest, "account id is required")
		return
	}

	id, err := strconv.Atoi(pathId)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err = req.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = d.service.Create(r.Context(), req, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.AccountNotFound):
			response.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.AccountIsNotActive):
			response.WriteError(w, http.StatusNotFound, err.Error())
		default:
			d.logger.Error().Err(err).Str("account_id", pathId).Msg("Error occured")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	d.logger.Info().Str("account_id", pathId).Str("amount", req.Amount.String()).Msg("Created deposit")
	response.WriteJson(w, http.StatusCreated, "created")
}
