package middleware

import (
	"net/http"

	"github.com/krishnasingh0210/blog/pkg/response"
)

// RequireRole restricts a route to one or more roles. Must run AFTER Auth,
// since it reads the role that Auth already put into context.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextRoleKey).(string)
			if !ok {
				response.Error(w, http.StatusForbidden, "role not found in context")
				return
			}
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Error(w, http.StatusForbidden, "insufficient permissions")
		})
	}
}