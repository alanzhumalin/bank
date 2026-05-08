package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/internal/handler"
	"github.com/alanzhumalin/bank/pkg/jwt"
)

type AuthMiddleware struct {
	TokenKey *string
}

func NewAuthMiddleWare(tokenKey *string) *AuthMiddleware {
	return &AuthMiddleware{
		TokenKey: tokenKey,
	}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			handler.WriteJson(w, http.StatusUnauthorized, "not authenticated")
			return
		}

		parts := strings.Split(header, " ")

		if len(parts) != 2 {
			handler.WriteJson(w, http.StatusUnauthorized, "incorrect authorization header")
			return
		}

		if parts[0] != "Bearer" {
			handler.WriteJson(w, http.StatusUnauthorized, "invalid auth type")
			return
		}

		claims, err := jwt.ParseAndValidateToken(parts[1], *a.TokenKey)

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrorNotValidToken):
				handler.WriteJson(w, http.StatusUnauthorized, err.Error())
			default:
				handler.WriteJson(w, http.StatusUnauthorized, "Not authorized")
			}

			return
		}

		ctx := context.WithValue(r.Context(), dto.UserKey{}, claims.UserId)
		ctx = context.WithValue(ctx, dto.RoleKey{}, claims.Role)
		ctx = context.WithValue(ctx, dto.SessionKey{}, claims.SessionId)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
