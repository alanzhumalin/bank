package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/middleware"
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

func TransactionRouter(th *transactionHandler, auth middleware.Middleware, ratelimit middleware.Middleware, rbac func(...string) middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{account_id}", middleware.Chain(http.HandlerFunc(th.GetByAccountId), auth, ratelimit))
	mux.Handle("GET /", middleware.Chain(http.HandlerFunc(th.GetByUserId), auth, ratelimit))

	return mux
}

func (th *transactionHandler) GetByUserId(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(dto.UserKey{}).(int)

	if userId <= 0 {
		response.WriteJson(w, http.StatusUnauthorized, "not authorized")
		return
	}

	query := r.URL.Query()

	var currencies *[]string

	if currencyRaw := query.Get("currencies"); currencyRaw != "" {

		cRaw := strings.Split(currencyRaw, ",")
		currencies = &cRaw
	}

	cursor := query.Get("cursor")
	limit := 20

	if limitRaw := query.Get("limit"); limitRaw != "" {
		limitNum, err := strconv.Atoi(limitRaw)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, "limit must be an integer")
			return
		}

		limit = limitNum
	}

	if limit <= 0 {
		response.WriteJson(w, http.StatusBadRequest, "limit must be non 0")
		return
	}

	if limit > 100 {
		limit = 100
	}

	res, err := th.service.GetByUserId(r.Context(), userId, cursor, limit, currencies)

	if err != nil {
		switch {

		default:
			th.logger.Error().Err(err).Msg("Error occured")
			response.WriteJson(w, http.StatusInternalServerError, "internal server error")

		}
		return
	}

	th.logger.Info().Str("user_id", strconv.Itoa(userId)).Msg("Get transactions by user_id")
	response.WriteJson(w, http.StatusOK, res)

}

func (th *transactionHandler) GetByAccountId(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("account_id")

	query := r.URL.Query()
	cursor := query.Get("cursor")

	limit := 20

	if limitRaw := query.Get("limit"); limitRaw != "" {
		parsedLimit, err := strconv.Atoi(limitRaw)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "limit must be an integer")
			return
		}
		limit = parsedLimit
	}

	if limit <= 0 {
		response.WriteError(w, http.StatusBadRequest, "limit must be not negative")
		return
	}

	if limit > 100 {
		limit = 100
	}

	userIdFromContext := r.Context().Value(dto.UserKey{}).(int)

	if pathId == "" {
		response.WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	if userIdFromContext <= 0 {
		response.WriteError(w, http.StatusBadRequest, "account_id is required")
		return
	}

	id, err := strconv.Atoi(pathId)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "account_id must be an integer")
		return
	}

	trs, err := th.service.GetByAccountId(r.Context(), id, limit, cursor, userIdFromContext)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorForBidden):
			response.WriteError(w, http.StatusForbidden, "forbidden")

		default:
			th.logger.Error().Err(err).Msg("Error in get transactions by account id")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
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
