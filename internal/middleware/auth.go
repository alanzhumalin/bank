package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/pkg/jwt"
	"github.com/alanzhumalin/bank/pkg/response"
)

type AuthMiddleware struct {
	TokenKey *string
}

func NewAuthMiddleWare(tokenKey *string) *AuthMiddleware {
	return &AuthMiddleware{
		TokenKey: tokenKey,
	}
}

func (a *AuthMiddleware) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				response.WriteJson(w, http.StatusForbidden, "forbidden")
				return
			}

			parts := strings.Split(header, " ")

			if len(parts) != 2 {
				response.WriteJson(w, http.StatusForbidden, "forbidden")
				return
			}

			if parts[0] != "Bearer" {
				response.WriteJson(w, http.StatusForbidden, "forbidden")
				return
			}

			claims, err := jwt.ParseAndValidateToken(parts[1], *a.TokenKey)

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrorNotValidToken):
					response.WriteJson(w, http.StatusUnauthorized, "unauthorized")

				default:
					response.WriteJson(w, http.StatusUnauthorized, "Not authorized")
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
}
