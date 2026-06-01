package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alanzhumalin/bank/internal/cache"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/alanzhumalin/bank/pkg/response"
)

var RateLimiterIsRequired = errors.New("Rate limiter is required")

type RateLimitMiddleware struct {
	rateLimiter *cache.RateLimiter
}

func NewRateLimiterMiddleware(rateLimiter *cache.RateLimiter) (RateLimitMiddleware, error) {
	if rateLimiter != nil {
		return RateLimitMiddleware{rateLimiter: rateLimiter}, nil
	}
	return RateLimitMiddleware{}, RateLimiterIsRequired
}

func (rm *RateLimitMiddleware) RateLimiterMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, ok := r.Context().Value(dto.UserKey{}).(int)
			if !ok {
				response.WriteError(w, http.StatusBadRequest, "user id is must be a number")
				return
			}

			if userId <= 0 {
				response.WriteError(w, http.StatusBadRequest, "user id is empty")
				return
			}

			id := strconv.Itoa(userId)

			key := fmt.Sprintf("rate:user_id:%s", id)

			check, err := rm.rateLimiter.IsAllowed(r.Context(), key)

			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			if !check {
				response.WriteError(w, http.StatusTooManyRequests, "too many requests")
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
