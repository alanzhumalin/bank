package handler

import (
	"net/http"

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

func DepositRouter(d *depositHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{id}/deposit", d.Create)
	return mux
}

func (d *depositHandler) Create(w http.ResponseWriter, r *http.Request) {

}
