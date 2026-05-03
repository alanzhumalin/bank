package handler

import "net/http"

func AccountRouter(accountHandler *accountHandler, depositHandler *depositHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", accountHandler.Create)
	mux.HandleFunc("DELETE /{account_id}", accountHandler.DeleteByID)
	mux.HandleFunc("GET /", accountHandler.GetAll)
	mux.HandleFunc("POST /{account_id}/deposit", depositHandler.Create)
	return mux
}
