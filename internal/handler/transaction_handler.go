package handler

import (
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/service"
	"github.com/alanzhumalin/bank/pkg/response"
	"github.com/rs/zerolog"
)

type transactionHandler struct {
	service service.TransactionService
	logger  zerolog.Logger
}

func NewTransactionHandler(service service.TransactionService, logger zerolog.Logger) *transactionHandler {
	return &transactionHandler{
		service: service,
		logger:  logger.With().Str("component", "transaction_handler").Logger(),
	}
}

func TransactionRouter(th *transactionHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", th.GetAll)
	mux.HandleFunc("GET /{account_id}", th.GetByAccountId)

	return mux
}

func (th *transactionHandler) GetByAccountId(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("account_id")

	if pathId == "" {
		response.WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	id, err := strconv.Atoi(pathId)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "account_id must be an integer")
		return
	}

	trs, err := th.service.GetByAccountId(r.Context(), id)

	if err != nil {
		th.logger.Error().Err(err).Msg("Error in get transactions by account id")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	th.logger.Info().Str("account_id", pathId).Msg("Get all transactions by account id")
	response.WriteJson(w, http.StatusOK, trs)

}

func (th *transactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	trs, err := th.service.GetAll(r.Context())

	if err != nil {
		th.logger.Error().Err(err).Msg("Error in get all transactions")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	th.logger.Info().Msg("Get all transactions")
	response.WriteJson(w, http.StatusOK, trs)

}
