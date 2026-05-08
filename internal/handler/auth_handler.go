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

type authHandler struct {
	userService service.UserService
	authService service.AuthService
	logger      zerolog.Logger
}

func AuthRouter(ah *authHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", ah.Register)
	mux.HandleFunc("POST /login", ah.Login)
	mux.HandleFunc("POST /refresh", ah.Refresh)
	mux.HandleFunc("POST /logout", ah.Logout)
	mux.HandleFunc("POST /logoutall", ah.LogoutFromAllDevices)

	return mux
}

func NewAuthHandler(userService service.UserService, authService service.AuthService, logger zerolog.Logger) *authHandler {
	return &authHandler{userService: userService, authService: authService, logger: logger.With().Str("component", "auth_handler").Logger()}
}

func (ah *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	metadata := GetRequestMetadata(r)

	token, err := ah.authService.Register(r.Context(), req, metadata.Ip, metadata.Agent)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorUserAlreadyExists):
			WriteError(w, http.StatusConflict, err.Error())
		default:
			ah.logger.Error().Err(err).Msg("error in registering user")
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	ah.logger.Info().Msg("Registered user")

	WriteJson(w, http.StatusCreated, token)

}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	metadata := GetRequestMetadata(r)

	token, err := ah.authService.Login(r.Context(), req, metadata.Ip, metadata.Agent)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrorUserNotFound):
			WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrorPasswordNotCorrect):
			WriteError(w, http.StatusNotAcceptable, err.Error())
		default:
			ah.logger.Error().Err(err).Msg("Error in login")
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Msg("Login")
	WriteJson(w, http.StatusOK, token)
}

func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionId, ok := r.Context().Value(dto.SessionKey{}).(string)

	if !ok {
		WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := ah.authService.Logout(r.Context(), sessionId); err != nil {
		switch {
		case errors.Is(err, domain.ErrorSessionNotFound):
			WriteError(w, http.StatusNotFound, "session not found")
		default:
			ah.logger.Error().Err(err).Msg("Error in logout")
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Str("session_id", sessionId).Msg("Logout")
	WriteJson(w, http.StatusOK, "")
}

func (ah *authHandler) LogoutFromAllDevices(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(dto.UserKey{}).(int)

	if !ok {
		WriteError(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := ah.authService.LogoutFromAllDevices(r.Context(), userId); err != nil {
		switch {
		case errors.Is(err, domain.ErrorSessionNotFound):
			WriteError(w, http.StatusNotFound, "sessions not found")
		default:
			ah.logger.Error().Err(err).Msg("Error in logout from all devices")
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Str("user_id", strconv.Itoa(userId)).Msg("Logout from all devices")
	WriteJson(w, http.StatusOK, "")
}

func (ah *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	sessionId, ok := r.Context().Value(dto.SessionKey{}).(string)
	if !ok {
		WriteError(w, http.StatusBadRequest, "bad request")
		return
	}
	userId, ok := r.Context().Value(dto.UserKey{}).(int)
	if !ok {
		WriteError(w, http.StatusBadRequest, "bad request")
		return

	}
	role, ok := r.Context().Value(dto.RoleKey{}).(string)

	if !ok {
		WriteError(w, http.StatusBadRequest, "bad request")
		return

	}

	token, err := ah.authService.UpdateSession(r.Context(), userId, role, sessionId)

	if err != nil {
		switch {
		case errors.Is(err, dto.SessionIdNotFound):
			WriteError(w, http.StatusBadRequest, "bad request")

		default:
			ah.logger.Error().Err(err).Msg("Error occured in updating the session")
			WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	ah.logger.Info().Str("session_id", sessionId).Msg("Update session")
	WriteJson(w, http.StatusOK, token)

}
