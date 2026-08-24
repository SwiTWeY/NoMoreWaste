package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/SwiTWeY/NoMoreWaste/api/auth"
	"github.com/SwiTWeY/NoMoreWaste/api/utils"
)

type cleContexte string

const CleClaims cleContexte = "claims"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entete := r.Header.Get("Authorization")
			if !strings.HasPrefix(entete, "Bearer ") {
				utils.Error(w, http.StatusUnauthorized, "token manquant")
				return
			}
			tokenStr := strings.TrimPrefix(entete, "Bearer ")
			claims, err := auth.ValiderToken(secret, tokenStr)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "token invalide")
				return
			}
			ctx := context.WithValue(r.Context(), CleClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
func Personnel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(CleClaims).(*auth.Claims)
		if !ok || !claims.EstPersonnel {
			utils.Error(w, http.StatusForbidden, "acces reserve au personnel")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ClaimsDepuis(r *http.Request) (*auth.Claims, bool) {
	c, ok := r.Context().Value(CleClaims).(*auth.Claims)
	return c, ok
}
