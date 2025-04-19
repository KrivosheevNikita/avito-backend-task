package middleware

import (
	"net/http"
	"strings"

	"pvz-service/internal/models"
	"pvz-service/pkg/handlerutils"
	"pvz-service/pkg/utils"

	"github.com/gorilla/mux"
)

var parseToken = utils.ParseToken

func RoleCheck(allowedRoles ...string) mux.MiddlewareFunc {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				handlerutils.WriteError(w, models.ErrForbidden)
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims, err := parseToken(tokenString)
			if err != nil {
				handlerutils.WriteError(w, models.ErrForbidden)
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				handlerutils.WriteError(w, models.ErrForbidden)
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				handlerutils.WriteError(w, models.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
