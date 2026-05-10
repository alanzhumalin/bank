package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/middleware"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/alanzhumalin/bank/pkg/response"
	"github.com/rs/zerolog"
)

func UserRouter(userHandler *userHandler, authMiddleware middleware.AuthMiddleware, rbac middleware.RbacMiddleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", middleware.Chain(http.HandlerFunc(userHandler.GetAll), authMiddleware.Middleware(), rbac.RBAC()))
	return mux
}

type userHandler struct {
	userService service.UserService
	log         zerolog.Logger
}

func NewUserHandler(service service.UserService, logger zerolog.Logger) *userHandler {
	return &userHandler{userService: service, log: logger.With().Str("component", "user_handler").Logger()}
}

func (u *userHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.GetAll(r.Context())

	if err != nil {
		u.log.Error().Err(err).Msg("Error in get all users")
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	u.log.Info().Msg("Get all users")

	response.WriteJson(w, http.StatusOK, users)
}

func (u *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req dto.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		u.log.Warn().Err(err).Msg("failed to decode request")
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = req.Validate()
	if err != nil {
		u.log.Warn().Err(err).Msg("Incorrect validation")
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, _, err = u.userService.Create(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorUserAlreadyExists):
			u.log.Warn().Err(err).Msg("User already exists")
			response.WriteError(w, http.StatusConflict, err.Error())
		default:
			u.log.Error().Err(err).Msg("failed to create user")

			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	u.log.Info().Str("Firstname", req.FirstName).Str("LastName", req.LastName).Str("PhoneNumber", req.PhoneNumber).Dur("duration", time.Since(start)).Msg("Created new user")

	response.WriteJson(w, http.StatusCreated, "created")

}

func (u *userHandler) GetByPhone(w http.ResponseWriter, r *http.Request) {

	phone := r.PathValue("phone")
	if phone == "" {
		u.log.Warn().Err(dto.ErrorPhoneNumRequired).Msg("Phone number is required")
		response.WriteError(w, http.StatusBadRequest, "Phone number is required")
		return
	}

	res, err := u.userService.GetByPhone(r.Context(), phone)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorUserNotFound):
			u.log.Warn().Err(err).Str("phone", phone).Msg("user not found")
			response.WriteError(w, http.StatusNotFound, err.Error())

		default:
			u.log.Error().Err(err).Str("phone", phone).Msg("Error get user by phone")
			response.WriteError(w, http.StatusInternalServerError, "interval server error")
		}
		return
	}

	u.log.Info().Str("phone_number", phone).Msg("Get user by phone number")

	response.WriteJson(w, http.StatusOK, res)
}
