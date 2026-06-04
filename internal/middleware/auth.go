package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alanzhumalin/bank/internal/cache"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/pkg/jwt"
	"github.com/alanzhumalin/bank/pkg/response"
)

type AuthMiddleware struct {
	TokenKey       *string
	TokenBlackList cache.TokenBlackList
}

func NewAuthMiddleWare(tokenKey *string, tokenBlackList cache.TokenBlackList) *AuthMiddleware {
	return &AuthMiddleware{
		TokenKey:       tokenKey,
		TokenBlackList: tokenBlackList,
	}
}

func (a *AuthMiddleware) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				response.WriteJson(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			parts := strings.Fields(header)

			if len(parts) != 2 {
				response.WriteJson(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if parts[0] != "Bearer" {
				response.WriteJson(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			claims, err := jwt.ParseAndValidateToken(parts[1], *a.TokenKey)

			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrorNotValidToken):
					response.WriteJson(w, http.StatusUnauthorized, "unauthorized")

				default:
					response.WriteJson(w, http.StatusUnauthorized, "unauthorized")
				}

				return
			}

			ok, err := a.TokenBlackList.Exists(r.Context(), claims.JTI)
			if err != nil {
				response.WriteJson(w, http.StatusInternalServerError, "internal server error")
				return
			}

			if ok {
				response.WriteJson(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), dto.UserKey{}, claims.UserId)
			ctx = context.WithValue(ctx, dto.RoleKey{}, claims.Role)
			ctx = context.WithValue(ctx, dto.SessionKey{}, claims.SessionId)
			ctx = context.WithValue(ctx, dto.JTIKey{}, claims.JTI)
			ctx = context.WithValue(ctx, dto.ExpKey{}, claims.ExpiresAt.Time)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
