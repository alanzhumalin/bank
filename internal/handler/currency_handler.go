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

func CurrencyRouter(c *currencyHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", c.Create)
	mux.HandleFunc("GET /", c.GetAll)
	mux.HandleFunc("DELETE /{id}", c.Delete)
	mux.HandleFunc("GET /{id}", c.GetById)

	return mux
}

type currencyHandler struct {
	service service.CurrencyService
	logger  zerolog.Logger
}

func NewCurrencyHandler(service service.CurrencyService, logger zerolog.Logger) *currencyHandler {
	return &currencyHandler{
		service: service,
		logger:  logger.With().Str("component", "currency_handler").Logger(),
	}
}

func (c *currencyHandler) GetById(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	if pathId == "" {
		response.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(pathId)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "id must be a number")
		return
	}

	res, err := c.service.GetById(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorCurrencyNotFound):
			response.WriteError(w, http.StatusBadRequest, "id must be a number")
		default:
			c.logger.Error().Err(err).Str("id", pathId).Msg("Error get by id currency_handler")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")

		}
		return
	}

	c.logger.Info().Str("id", pathId).Msg("Get by id ")
	response.WriteJson(w, http.StatusOK, res)

}

func (c *currencyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNewCurrencyRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = c.service.Create(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, domain.ErrorCurrencyAlreadyExists):
			c.logger.Warn().Str("code", req.Code).Msg("currency already exists")
			response.WriteError(w, http.StatusConflict, err.Error())
		default:
			c.logger.Error().Err(err).Str("name", req.Name).Str("code", req.Code).Str("symbol", req.Symbol).Msg("Error create new currency")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.logger.Info().Str("code", req.Code).Msg("Created currency")

	response.WriteJson(w, http.StatusCreated, "created")

}

func (c *currencyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")

	if pathId == "" {
		response.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(pathId)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "id must be number")
		return
	}

	err = c.service.Delete(r.Context(), id)

	if err != nil {
		c.logger.Error().Err(err).Msg("Error in deleting currency")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	c.logger.Info().Str("id", pathId).Msg("Delete currency")
	response.WriteJson(w, http.StatusNoContent, "deleted")

}

func (c *currencyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.GetAll(r.Context())

	if err != nil {
		c.logger.Error().Err(err).Msg("Error in get all currencies")
		response.WriteJson(w, http.StatusInternalServerError, "internal server error")
		return
	}

	c.logger.Info().Msg("Get all currencies")
	response.WriteJson(w, http.StatusOK, res)

}
