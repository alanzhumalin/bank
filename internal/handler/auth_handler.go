package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/middleware"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/alanzhumalin/bank/pkg/response"
	"github.com/rs/zerolog"
)

type authHandler struct {
	userService service.UserService
	authService service.AuthService
	logger      zerolog.Logger
}

func AuthRouter(ah *authHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /register", http.HandlerFunc(ah.Register))
	mux.Handle("POST /login", http.HandlerFunc(ah.Login))
	mux.Handle("POST /refresh", middleware.Chain(http.HandlerFunc(ah.Refresh)))
	mux.Handle("POST /logout", middleware.Chain(http.HandlerFunc(ah.Logout), auth))
	mux.Handle("POST /logoutall", middleware.Chain(http.HandlerFunc(ah.LogoutFromAllDevices), auth))

	return mux
}

func NewAuthHandler(userService service.UserService, authService service.AuthService, logger zerolog.Logger) *authHandler {
	return &authHandler{userService: userService, authService: authService, logger: logger.With().Str("component", "auth_handler").Logger()}
}

func (ah *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	metadata := GetRequestMetadata(r)

	token, err := ah.authService.Register(r.Context(), req, metadata.Ip, metadata.Agent)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorUserAlreadyExists):
			response.WriteError(w, http.StatusConflict, err.Error())
		default:
			ah.logger.Error().Err(err).Msg("error in registering user")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	ah.logger.Info().Msg("Registered user")

	response.WriteJson(w, http.StatusCreated, token)

}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	metadata := GetRequestMetadata(r)

	token, err := ah.authService.Login(r.Context(), req, metadata.Ip, metadata.Agent)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorUserNotFound):
			response.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrorPasswordNotCorrect):
			response.WriteError(w, http.StatusNotAcceptable, err.Error())
		default:
			ah.logger.Error().Err(err).Msg("Error in login")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Msg("Login")
	response.WriteJson(w, http.StatusOK, token)
}

func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionId, ok := r.Context().Value(dto.SessionKey{}).(string)

	if !ok {
		response.WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	jti, ok := r.Context().Value(dto.JTIKey{}).(string)

	if !ok {
		response.WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	exp, ok := r.Context().Value(dto.ExpKey{}).(time.Time)

	if !ok {
		response.WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := ah.authService.Logout(r.Context(), sessionId, jti, exp); err != nil {
		switch {
		case errors.Is(err, domain.ErrorSessionNotFound):
			response.WriteError(w, http.StatusNotFound, "session not found")
		default:
			ah.logger.Error().Err(err).Msg("Error in logout")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Str("session_id", sessionId).Msg("Logout")
	response.WriteJson(w, http.StatusOK, "logout")
}

func (ah *authHandler) LogoutFromAllDevices(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(dto.UserKey{}).(int)

	if !ok {
		response.WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := ah.authService.LogoutFromAllDevices(r.Context(), userId); err != nil {
		switch {
		case errors.Is(err, domain.ErrorSessionNotFound):
			response.WriteError(w, http.StatusNotFound, "sessions not found")
		default:
			ah.logger.Error().Err(err).Msg("Error in logout from all devices")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Str("user_id", strconv.Itoa(userId)).Msg("Logout from all devices")
	response.WriteJson(w, http.StatusOK, "")
}

func (ah *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	var req dto.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	token, sessionId, err := ah.authService.UpdateSession(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, dto.SessionIdNotFound):
			response.WriteError(w, http.StatusNotFound, "Session not found")

		default:
			ah.logger.Error().Err(err).Msg("Error occured in updating the session")
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Str("session_id", sessionId).Msg("Update session")
	response.WriteJson(w, http.StatusOK, token)

}
