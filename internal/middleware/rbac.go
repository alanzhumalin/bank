package middleware

import (
	"net/http"
	"slices"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/pkg/response"
)

type RbacMiddleware struct{}

func NewRbacMiddleware() *RbacMiddleware {
	return &RbacMiddleware{}
}

func (m *RbacMiddleware) RBAC(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(dto.RoleKey{}).(string)
			if !ok {
				response.WriteError(w, http.StatusForbidden, "forbidden")
				return
			}

			if ok := slices.Contains(roles, role); ok {
				next.ServeHTTP(w, r)
				return
			}

			response.WriteError(w, http.StatusForbidden, "forbidden")
		})
	}
}
