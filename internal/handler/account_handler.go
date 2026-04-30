package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/rs/zerolog"
)

type accountHandler struct {
	service service.AccountService
	logger  zerolog.Logger
}

func NewAccountHandler(service service.AccountService, logger zerolog.Logger) *accountHandler {
	return &accountHandler{
		service: service,
		logger:  logger.With().Str("component", "account_handler").Logger(),
	}
}

func AccountRouter(a *accountHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", a.Create)

	return mux
}

func (a *accountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.service.Create(r.Context(), req)

	if err != nil {
		a.logger.Error().Err(err).Msg("Error in creating the account")
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.logger.Info().Msg("Created a new account")
	WriteJson(w, http.StatusCreated, "created")
}

func (a *accountHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")

	if pathId == "" {
		WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(pathId)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if err = a.service.DeleteById(r.Context(), id); err != nil {
		a.logger.Error().Err(err).Msg("Error in deleting account by id")
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.logger.Info().Str("id", pathId).Msg("Deleted account by id")
	WriteJson(w, http.StatusNoContent, "deleted")

}

func (a *accountHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	accounts, err := a.service.GetAll(r.Context())

	if err != nil {
		a.logger.Error().Err(err).Msg("Error in getting all accounts")
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.logger.Info().Msg("get all accounts")
	WriteJson(w, http.StatusOK, accounts)
}
