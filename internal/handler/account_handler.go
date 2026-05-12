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

func AccountRouter(accountHandler *accountHandler, authMiddleware middleware.Middleware, rbac func(...string) middleware.Middleware) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /", middleware.Chain(http.HandlerFunc(accountHandler.Create), authMiddleware))
	mux.Handle("DELETE /{account_id}", middleware.Chain(http.HandlerFunc(accountHandler.DeleteByID), authMiddleware))
	mux.Handle("GET /", middleware.Chain(http.HandlerFunc(accountHandler.GetAll), authMiddleware, rbac("admin")))
	return mux
}

func (a *accountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.service.Create(r.Context(), req)

	if err != nil {
		a.logger.Error().Err(err).Msg("Error in creating the account")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.logger.Info().Msg("Created a new account")
	response.WriteJson(w, http.StatusCreated, "created")
}

func (a *accountHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("account_id")

	if pathId == "" {
		response.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(pathId)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if err = a.service.DeleteById(r.Context(), id); err != nil {
		a.logger.Error().Err(err).Msg("Error in deleting account by id")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.logger.Info().Str("id", pathId).Msg("Deleted account by id")
	response.WriteJson(w, http.StatusNoContent, "deleted")

}

func (a *accountHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	accounts, err := a.service.GetAll(r.Context())

	if err != nil {
		a.logger.Error().Err(err).Msg("Error in getting all accounts")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	a.logger.Info().Msg("get all accounts")
	response.WriteJson(w, http.StatusOK, accounts)
}
